package server

import (
	"context"
	pbcommon "github.com/kic/feed/pkg/proto/common"
)

type FriendServicer interface {
	GetFriendsForUser(context.Context, int64, string) ([]uint64, error)
	GetConnectionBetweenUsers(context.Context, int64, int64, string) (float32, error)
}

type UserServicer interface {
	GetUserNameForID(ctx context.Context, userID int64, authCredentials string) (string, error)
}

type MediaServicer interface {
	GetFilesForUser(ctx context.Context, userID int64, authCredentials string) ([]*pbcommon.File, error)
}

type HealthServicer interface {
	GetMentalHealthScoreForUser(ctx context.Context, userID int64, authCredentials string) (int32, error)
}
