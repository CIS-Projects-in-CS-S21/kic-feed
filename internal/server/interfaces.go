package server

import "context"

type FriendServicer interface {
	GetFriendsForUser(context.Context, int64, string) ([]uint64, error)
	GetConnectionBetweenUsers(context.Context, int64, int64, string) (float32, error)
}

type UserServicer interface {
}

type MediaServicer interface {
}
