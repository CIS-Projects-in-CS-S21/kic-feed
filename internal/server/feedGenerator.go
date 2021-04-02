package server

import (
	"context"
	pbcommon "github.com/kic/feed/pkg/proto/common"
	"go.uber.org/zap"
)

type FeedGenerator struct {
	logger *zap.SugaredLogger

	friendsClient FriendServicer
	usersClient   UserServicer
	mediaClient   MediaServicer
}

func NewFeedGenerator(
	logger *zap.SugaredLogger,
	friendsClient FriendServicer,
	usersClient UserServicer,
	mediaClient MediaServicer,
) *FeedGenerator {
	return &FeedGenerator{
		logger:        logger,
		friendsClient: friendsClient,
		usersClient:   usersClient,
		mediaClient:   mediaClient,
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
	return nil, nil
}

func (f *FeedGenerator) getUserNameForID(ctx context.Context, userID int64, authCredentials string) (string, error) {
	return "", nil
}

func (f *FeedGenerator) rankAndSortPosts(
	ctx context.Context,
	userID int64,
	authCredentials string,
	posts []*pbcommon.File,
) error {
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
			return nil, err
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

	return allPosts, nil
}
