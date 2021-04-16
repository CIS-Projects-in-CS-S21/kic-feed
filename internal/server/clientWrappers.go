package server

import (
	"context"
	"go.uber.org/zap"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pbcommon "github.com/kic/feed/pkg/proto/common"
	pbfriends "github.com/kic/feed/pkg/proto/friends"
	pbhealth "github.com/kic/feed/pkg/proto/health"
	pbmedia "github.com/kic/feed/pkg/proto/media"
	pbusers "github.com/kic/feed/pkg/proto/users"
)

type FriendClientWrapper struct {
	friendsClient pbfriends.FriendsClient
	logger        *zap.SugaredLogger
}

func NewFriendClientWrapper(conn *grpc.ClientConn, logger *zap.SugaredLogger) *FriendClientWrapper {
	return &FriendClientWrapper{
		friendsClient: pbfriends.NewFriendsClient(conn),
		logger:        logger,
	}
}

func (f *FriendClientWrapper) GetFriendsForUser(
	ctx context.Context,
	userID int64,
	authCredentials string,
) ([]uint64, error) {
	f.logger.Debugf("Getting friends for %v", userID)
	md := metadata.Pairs("Authorization", authCredentials)
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pbfriends.GetFriendsForUserRequest{
		User: &pbcommon.User{
			UserID: userID,
		},
	}
	resp, err := f.friendsClient.GetFriendsForUser(ctx, req)

	if err != nil {
		f.logger.Debugf("Failed to get friends for %v, returning %v", userID, err)
		return nil, err
	}

	f.logger.Debugf("Got friends for %v, they are %v", userID, resp.Friends)

	return resp.Friends, nil
}

func (f *FriendClientWrapper) GetConnectionBetweenUsers(
	ctx context.Context,
	uid1 int64,
	uid2 int64,
	authCredentials string,
) (float32, error) {
	md := metadata.Pairs("Authorization", authCredentials)
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pbfriends.GetConnectionBetweenUsersRequest{
		FirstUserID:  uint64(uid1),
		SecondUserID: uint64(uid2),
	}

	resp, err := f.friendsClient.GetConnectionBetweenUsers(ctx, req)

	if err != nil {
		f.logger.Debugf("Failed to get friend connection for %v and %v, returning %v", uid1, uid2, err)
		return 0.0, err
	}

	f.logger.Debugf("Got connection for %v and %v, it is %v", uid1, uid2, resp.ConnectionStrength)

	return resp.ConnectionStrength, nil
}

type MediaClientWrapper struct {
	mediaClient pbmedia.MediaStorageClient
	logger      *zap.SugaredLogger
}

func NewMediaClientWrapper(conn *grpc.ClientConn, logger *zap.SugaredLogger) *MediaClientWrapper {
	return &MediaClientWrapper{
		mediaClient: pbmedia.NewMediaStorageClient(conn),
		logger:      logger,
	}
}

func (m *MediaClientWrapper) GetFilesForUser(
	ctx context.Context,
	userID int64,
	authCredentials string,
) ([]*pbcommon.File, error) {
	md := metadata.Pairs("Authorization", authCredentials)
	ctx = metadata.NewOutgoingContext(ctx, md)

	fileMetadata := make(map[string]string)

	fileMetadata["userID"] = strconv.FormatInt(userID, 10)

	req := &pbmedia.GetFilesByMetadataRequest{
		DesiredMetadata: fileMetadata,
		Strictness:      pbmedia.MetadataStrictness_STRICT,
	}

	resp, err := m.mediaClient.GetFilesWithMetadata(ctx, req)

	if err != nil {
		m.logger.Debugf("Failed to get files for %v, returning %v", userID, err)
		return nil, err
	}

	m.logger.Debugf("Got files for %v, returning %v", userID, resp.FileInfos)

	return resp.FileInfos, nil
}

type UserClientWrapper struct {
	usersClient pbusers.UsersClient
	logger      *zap.SugaredLogger
}

func NewUserClientWrapper(conn *grpc.ClientConn, logger *zap.SugaredLogger) *UserClientWrapper {
	return &UserClientWrapper{
		usersClient: pbusers.NewUsersClient(conn),
		logger:      logger,
	}
}

func (w *UserClientWrapper) GetUserNameForID(
	ctx context.Context,
	userID int64,
	authCredentials string,
) (string, error) {
	md := metadata.Pairs("Authorization", authCredentials)
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pbusers.GetUserNameByIDRequest{UserID: userID}

	resp, err := w.usersClient.GetUserNameByID(ctx, req)

	if err != nil {
		w.logger.Debugf("Failed to get username for %v, returning %v", userID, err)
		return "", err
	}

	w.logger.Debugf("Got username for %v, returning %v", userID, resp.Username)

	return resp.Username, nil
}

type HealthClientWrapper struct {
	healthClient pbhealth.HealthTrackingClient
	logger       *zap.SugaredLogger
}

func NewHealthClientWrapper(conn *grpc.ClientConn, logger *zap.SugaredLogger) *HealthClientWrapper {
	return &HealthClientWrapper{
		healthClient: pbhealth.NewHealthTrackingClient(conn),
		logger:       logger,
	}
}

func (h *HealthClientWrapper) GetMentalHealthScoreForUser(
	ctx context.Context,
	userID int64,
	authCredentials string,
) (int32, error) {

	md := metadata.Pairs("Authorization", authCredentials)
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pbhealth.GetMentalHealthScoreForUserRequest{UserID: userID}

	resp, err := h.healthClient.GetMentalHealthScoreForUser(ctx, req)

	if err != nil {
		h.logger.Debugf("Failed to get health score for %v, returning %v", userID, err)
		return 0, err
	}

	h.logger.Debugf("Got health score for %v, returning %v", userID, resp.Score)

	return resp.GetScore(), nil
}
