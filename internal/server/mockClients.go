package server

import (
	"context"

	pbcommon "github.com/kic/feed/pkg/proto/common"
)

type MockFriendClient struct {
	friendsList map[int64][]uint64
	connections map[int64][]float32
}

func NewMockFriendClient(friendsList map[int64][]uint64, connections map[int64][]float32) *MockFriendClient {
	return &MockFriendClient{
		friendsList: friendsList,
		connections: connections,
	}
}

func (f *MockFriendClient) GetFriendsForUser(
	ctx context.Context,
	userID int64,
	authCredentials string,
) ([]uint64, error) {
	return f.friendsList[userID], nil
}

func (f *MockFriendClient) GetConnectionBetweenUsers(
	ctx context.Context,
	uid1 int64,
	uid2 int64,
	authCredentials string,
) (float32, error) {
	for i := 0; i < len(f.friendsList[uid1]); i++ {
		if f.friendsList[uid1][i] == uint64(uid2) {
			return f.connections[uid1][i], nil
		}
	}

	return 0.0, nil
}

type MockMediaClient struct {
	files map[int64][]*pbcommon.File
}

func NewMockMediaClient(files map[int64][]*pbcommon.File) *MockMediaClient {
	return &MockMediaClient{
		files: files,
	}
}

func (m *MockMediaClient) GetFilesForUser(
	ctx context.Context,
	userID int64,
	authCredentials string,
) ([]*pbcommon.File, error) {
	return m.files[userID], nil
}

type MockUserClient struct {
	usernames map[int64]string
}

func NewMockUserClient(usernames map[int64]string) *MockUserClient {
	return &MockUserClient{
		usernames: usernames,
	}
}

func (w *MockUserClient) GetUserNameForID(
	ctx context.Context,
	userID int64,
	authCredentials string,
) (string, error) {
	return w.usernames[userID], nil
}

type MockHealthClient struct {
	scores map[int64]int32
}

func NewMockHealthClient(scores map[int64]int32) *MockHealthClient {
	return &MockHealthClient{scores: scores}
}

func (h *MockHealthClient) GetMentalHealthScoreForUser(
	ctx context.Context,
	userID int64,
	authCredentials string,
) (int32, error) {
	return h.scores[userID], nil
}
