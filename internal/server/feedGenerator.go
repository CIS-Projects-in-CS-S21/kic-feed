package server

import (
	"context"
	"errors"
	pbcommon "github.com/kic/feed/pkg/proto/common"
	"go.uber.org/zap"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

const (
	healthThreshold = -10
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
	// fisher-yates shuffle the array to try and break up posts by a given user
	for i := range posts {
		j := rand.Intn(i + 1)
		posts[i], posts[j] = posts[j], posts[i]
	}

	// we want chronological order to be the principle ordering, so the first index will be the most recent post
	sort.SliceStable(posts, func(i, j int) bool {
		if posts[i].DateStored.Year > posts[j].DateStored.Year {
			return true
		}
		if posts[i].DateStored.Month > posts[j].DateStored.Month {
			return true
		}
		if posts[i].DateStored.Day > posts[j].DateStored.Day {
			return true
		}
		return false
	})

	// since we do not distinguish between times in a given day, we sort today further by friend strength
	dateToday := time.Now()

	day := dateToday.Day()

	endFriendIndex := 0

	for {
		if int(posts[endFriendIndex].DateStored.Day) != day {
			break
		}
		endFriendIndex += 1
	}

	// we have multiple posts for today, do the sort
	if endFriendIndex > 2 {
		f.logger.Debug("Sorting by connection strength")
		todaySlice := posts[:endFriendIndex+1]
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

func (f *FeedGenerator) injectMentalHealthPosts(
	ctx context.Context,
	userID int64,
	authCredentials string,
	posts []*pbcommon.File,
) error {
	score, err := f.healthClient.GetMentalHealthScoreForUser(ctx, userID, authCredentials)

	if err != nil {
		f.logger.Errorf("Failed to get health score for user: %v", err)
		return err
	}

	if score < healthThreshold {
		// inject posts here
	}

	return nil
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

	// Now we rank the posts with the algorithm
	err = f.rankAndSortPosts(ctx, userID, authCredentials, allPosts)

	if err != nil {
		f.logger.Debugf("Failed to rank files for uid %v, err: %v", userID, err)
		return nil, err
	}

	// Finally inject mental health posts if needed

	err = f.injectMentalHealthPosts(ctx, userID, authCredentials, allPosts)

	if err != nil {
		f.logger.Debugf("Failed to inject posts for uid %v, err: %v", userID, err)
		return nil, err
	}

	return allPosts, nil
}
