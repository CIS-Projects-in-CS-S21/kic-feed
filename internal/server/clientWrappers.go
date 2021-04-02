package server

import (
	"context"
	"fmt"
	pbcommon "github.com/kic/feed/pkg/proto/common"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	pbfriends "github.com/kic/feed/pkg/proto/friends"
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
	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", authCredentials))
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

func (f *FriendClientWrapper) GetConnectionBetweenUsers(ctx context.Context, uid1 int64, uid2 int64, authCredentials string) (float32, error) {
	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", authCredentials))
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

type UserClientWrapper struct {
	usersClient pbusers.UsersClient
}

func NewUserClientWrapper(conn *grpc.ClientConn) *UserClientWrapper {
	return &UserClientWrapper{
		usersClient: pbusers.NewUsersClient(conn),
	}
}
