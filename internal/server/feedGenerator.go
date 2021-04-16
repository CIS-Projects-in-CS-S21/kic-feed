package server

import (
	"context"
	"errors"
	pbcommon "github.com/kic/feed/pkg/proto/common"
	"go.uber.org/zap"
	"math/rand"
	"sort"
	"strconv"
)

const (
	healthThreshold = 0
)

type FeedGenerator struct {
	logger *zap.SugaredLogger

	friendsClient FriendServicer
	usersClient   UserServicer
	mediaClient   MediaServicer
	healthClient  HealthServicer
}

func NewFeedGenerator(
	logger *zap.SugaredLogger,
	friendsClient FriendServicer,
	usersClient UserServicer,
	mediaClient MediaServicer,
	healthClient HealthServicer,
) *FeedGenerator {
	return &FeedGenerator{
		logger:        logger,
		friendsClient: friendsClient,
		usersClient:   usersClient,
		mediaClient:   mediaClient,
		healthClient:  healthClient,
	}
}

func (f *FeedGenerator) getFriendsForUser(ctx context.Context, userID int64, authCredentials string) ([]uint64, error) {
	friends, err := f.friendsClient.GetFriendsForUser(ctx, userID, authCredentials)

	if err != nil {
		f.logger.Errorf("Failed to get friends for user: %v", err)
		return nil, err
	}

	return friends, nil
}

func (f *FeedGenerator) getFilesForFriend(
	ctx context.Context,
	userID int64,
	authCredentials string,
) ([]*pbcommon.File, error) {
	files, err := f.mediaClient.GetFilesForUser(ctx, userID, authCredentials)

	if err != nil {
		f.logger.Errorf("Failed to get files for user: %v", err)
		return nil, err
	}

	return files, nil
}

func (f *FeedGenerator) getUserNameForID(ctx context.Context, userID int64, authCredentials string) (string, error) {
	username, err := f.usersClient.GetUserNameForID(ctx, userID, authCredentials)

	if err != nil {
		f.logger.Errorf("Failed to get username for user: %v", err)
		return "", err
	}

	if username == "" {
		f.logger.Errorf("Username blank for user: %v", err)
		return "", errors.New("empty username for friend")
	}

	return username, nil
}

func (f *FeedGenerator) rankAndSortPosts(
	ctx context.Context,
	userID int64,
	authCredentials string,
	posts []*pbcommon.File,
) error {
	// fisher-yates shuffle the array to try and break up posts by a given user, since we are using
	// stable sorts later on
	for i := range posts {
		j := rand.Intn(i + 1)
		posts[i], posts[j] = posts[j], posts[i]
	}

	// we want chronological order to be the principle ordering, so the first index will be the most recent post
	sort.SliceStable(posts, func(i, j int) bool {
		if posts[i].DateStored.Year < posts[j].DateStored.Year {
			return false
		} else if posts[i].DateStored.Year > posts[j].DateStored.Year {
			return true
		}

		if posts[i].DateStored.Month > posts[j].DateStored.Month {
			return true
		} else if posts[i].DateStored.Month < posts[j].DateStored.Month {
			return false
		} else if posts[i].DateStored.Day > posts[j].DateStored.Day {
			return true
		} else if posts[i].DateStored.Day < posts[j].DateStored.Day {
			return false
		}

		return true
	})

	// since we do not distinguish between times in a given day, we sort the latest day further by friend strength
	endFriendIndex := 0
	day := posts[endFriendIndex].DateStored.Day
	month := posts[endFriendIndex].DateStored.Month

	for endFriendIndex < len(posts) {
		if posts[endFriendIndex].DateStored.Day != day || posts[endFriendIndex].DateStored.Month != month {
			break
		}
		endFriendIndex += 1
	}

	// we have multiple posts for today, do the sort
	if endFriendIndex > 2 {
		f.logger.Debug("Sorting by connection strength")
		todaySlice := posts[:endFriendIndex]
		strengthMap := make(map[int]float32)

		for i := 0; i < endFriendIndex; i++ {
			uid, err := strconv.Atoi(todaySlice[i].Metadata["userID"])
			if err != nil {
				return err
			}
			if _, ok := strengthMap[uid]; !ok {
				strength, err := f.friendsClient.GetConnectionBetweenUsers(ctx, userID, int64(uid), authCredentials)
				if err != nil {
					return err
				}
				strengthMap[uid] = strength
			}
		}

		sort.SliceStable(todaySlice, func(i, j int) bool {
			user1ID, _ := strconv.Atoi(todaySlice[i].Metadata["userID"])
			user2ID, _ := strconv.Atoi(todaySlice[j].Metadata["userID"])
			return strengthMap[user1ID] < strengthMap[user2ID]
		})
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (f *FeedGenerator) injectMentalHealthPosts(
	ctx context.Context,
	userID int64,
	authCredentials string,
	posts []*pbcommon.File,
) ([]*pbcommon.File, error) {
	score, err := f.healthClient.GetMentalHealthScoreForUser(ctx, userID, authCredentials)

	if err != nil {
		f.logger.Errorf("Failed to get health score for user: %v", err)
		return posts, err
	}

	if score < healthThreshold {
		healthPosts, err := f.mediaClient.GetFilesForUser(ctx, 150, authCredentials)
		if err != nil {
			f.logger.Errorf("Failed to get mental health posts for user: %v", err)
			return posts, err
		}
		// inject posts here
		numInject := min(len(healthPosts), (len(posts)/5) + 1)
		healthPostIndex := 0
		returnSize := numInject + len(posts)
		f.logger.Debugf("Injecting %v mental health posts", numInject)
		toReturn := make([]*pbcommon.File, returnSize)

		for i := 0; i < returnSize; i++ {
			// every 5th post we make a mental health post
			if i%5 == 0 {
				toReturn[i] = healthPosts[healthPostIndex]
				healthPostIndex++
			} else {
				toReturn[i] = posts[i-healthPostIndex]
			}
		}
		return toReturn, nil
	}

	return posts, nil
}

func (f *FeedGenerator) GenerateFeedForUser(
	ctx context.Context,
	userID int64,
	authCredentials string,
) ([]*pbcommon.File, error) {
	// first we fetch all friends of the target user
	friendIDs, err := f.getFriendsForUser(ctx, userID, authCredentials)
	if err != nil {
		return nil, err
	}

	allPosts := make([]*pbcommon.File, 0)

	// then we fetch all the usernames of each friend along with their posts,
	// and attribute a username with a post. All posts are appended to one large slice
	for _, friendID := range friendIDs {
		userName, err := f.getUserNameForID(ctx, int64(friendID), authCredentials)
		if err != nil {
			f.logger.Debugf("Failed to get username for uid %v, err: %v", friendID, err)
			continue
		}

		files, err := f.getFilesForFriend(ctx, int64(friendID), authCredentials)

		if err != nil {
			f.logger.Debugf("Failed to get files for uid %v, err: %v", friendID, err)
			return nil, err
		}

		for _, file := range files {
			file.Metadata["posterUsername"] = userName
		}

		allPosts = append(allPosts, files...)
	}

	if len(allPosts) < 2 {
		return allPosts, nil
	}

	// Now we rank the posts with the algorithm
	err = f.rankAndSortPosts(ctx, userID, authCredentials, allPosts)

	if err != nil {
		f.logger.Debugf("Failed to rank files for uid %v, err: %v", userID, err)
		return nil, err
	}

	// Finally inject mental health posts if needed

	allPosts, err = f.injectMentalHealthPosts(ctx, userID, authCredentials, allPosts)

	if err != nil {
		f.logger.Debugf("Failed to inject posts for uid %v, err: %v", userID, err)
		return nil, err
	}

	return allPosts, nil
}
