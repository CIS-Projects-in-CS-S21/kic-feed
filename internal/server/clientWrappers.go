package server

import (
	"context"
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
}

func NewFriendClientWrapper(conn *grpc.ClientConn) *FriendClientWrapper {
	return &FriendClientWrapper{
		friendsClient: pbfriends.NewFriendsClient(conn),
	}
}

func (f *FriendClientWrapper) GetFriendsForUser(
	ctx context.Context,
	userID int64,
	authCredentials string,
) ([]uint64, error) {
	md := metadata.Pairs("Authorization", authCredentials)
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pbfriends.GetFriendsForUserRequest{
		User: &pbcommon.User{
			UserID: userID,
		},
	}
	resp, err := f.friendsClient.GetFriendsForUser(ctx, req)

	if err != nil {
		return nil, err
	}

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
		return 0.0, err
	}

	return resp.ConnectionStrength, nil
}

type MediaClientWrapper struct {
	mediaClient pbmedia.MediaStorageClient
}

func NewMediaClientWrapper(conn *grpc.ClientConn) *MediaClientWrapper {
	return &MediaClientWrapper{
		mediaClient: pbmedia.NewMediaStorageClient(conn),
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
		return nil, err
	}

	return resp.FileInfos, nil
}

type UserClientWrapper struct {
	usersClient pbusers.UsersClient
}

func NewUserClientWrapper(conn *grpc.ClientConn) *UserClientWrapper {
	return &UserClientWrapper{
		usersClient: pbusers.NewUsersClient(conn),
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
		return "", err
	}

	return resp.Username, nil
}

type HealthClientWrapper struct {
	healthClient pbhealth.HealthTrackingClient
}

func NewHealthClientWrapper(conn *grpc.ClientConn) *HealthClientWrapper {
	return &HealthClientWrapper{
		healthClient: pbhealth.NewHealthTrackingClient(conn),
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
		return 0, err
	}

	return resp.Score, nil
}
