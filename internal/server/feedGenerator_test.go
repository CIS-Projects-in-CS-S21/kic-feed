package server

import (
	"context"
	"fmt"
	pbcommon "github.com/kic/feed/pkg/proto/common"
	"testing"

	"go.uber.org/zap"

	"github.com/kic/feed/pkg/logging"
)

var gen *FeedGenerator

func configureFriends() FriendServicer {

	friendsList := map[int64][]uint64{
		1: {2, 3, 5},
		2: {1, 4, 5},
		3: {1, 4},
		4: {2, 3, 5},
		5: {1, 2, 4},
	}
	connections := map[int64][]float32{
		1: {1.0, 1.5, 0.5},
		2: {1.0, 1.0, 2.0},
		3: {1.5, 1.5},
		4: {1.0, 1.5, 0.5},
		5: {0.5, 2.0, 0.5},
	}

	return NewMockFriendClient(friendsList, connections)
}

func configureUsers() UserServicer {

	usernames := map[int64]string{
		1: "user1",
		2: "user2",
		3: "user3",
		4: "user4",
		5: "user5",
	}

	return NewMockUserClient(usernames)
}

func configureMedia() MediaServicer {

	files := map[int64][]*pbcommon.File{
		2: {
			&pbcommon.File{
				FileName:     "u2post1",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "2",
				},
				DateStored: &pbcommon.Date{
					Year:  2021,
					Month: 2,
					Day:   2,
				},
			},
			&pbcommon.File{
				FileName:     "u2post2",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "2",
				},
				DateStored: &pbcommon.Date{
					Year:  2021,
					Month: 1,
					Day:   2,
				},
			},
			&pbcommon.File{
				FileName:     "u2post3",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "2",
				},
				DateStored: &pbcommon.Date{
					Year:  2020,
					Month: 1,
					Day:   1,
				},
			},
		},
		3: {
			&pbcommon.File{
				FileName:     "u3post1",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "3",
				},
				DateStored: &pbcommon.Date{
					Year:  2021,
					Month: 2,
					Day:   2,
				},
			},
			&pbcommon.File{
				FileName:     "u3post2",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "3",
				},
				DateStored: &pbcommon.Date{
					Year:  2021,
					Month: 1,
					Day:   2,
				},
			},
			&pbcommon.File{
				FileName:     "u3post3",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "3",
				},
				DateStored: &pbcommon.Date{
					Year:  2020,
					Month: 1,
					Day:   1,
				},
			},
		},
		5: {
			&pbcommon.File{
				FileName:     "u5post1",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "5",
				},
				DateStored: &pbcommon.Date{
					Year:  2021,
					Month: 2,
					Day:   2,
				},
			},
			&pbcommon.File{
				FileName:     "u5post2",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "5",
				},
				DateStored: &pbcommon.Date{
					Year:  2021,
					Month: 1,
					Day:   2,
				},
			},
			&pbcommon.File{
				FileName:     "u5post3",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "5",
				},
				DateStored: &pbcommon.Date{
					Year:  2020,
					Month: 3,
					Day:   1,
				},
			},
		},
		-1: {
			&pbcommon.File{
				FileName:     "mentalHealth1",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "-1",
				},
				DateStored: &pbcommon.Date{
					Year:  2021,
					Month: 2,
					Day:   2,
				},
			},
			&pbcommon.File{
				FileName:     "mentalHealth2",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "-1",
				},
				DateStored: &pbcommon.Date{
					Year:  2021,
					Month: 1,
					Day:   2,
				},
			},
			&pbcommon.File{
				FileName:     "mentalHealth3",
				FileLocation: "",
				Metadata: map[string]string{
					"userID": "-1",
				},
				DateStored: &pbcommon.Date{
					Year:  2020,
					Month: 3,
					Day:   1,
				},
			},
		},
	}

	return NewMockMediaClient(files)
}

func configureHealth() HealthServicer {

	health := map[int64]int32{
		1: 2,
		2: -12,
		3: 1,
		4: -20,
		5: 5,
	}

	return NewMockHealthClient(health)
}

func TestMain(m *testing.M) {
	gen = NewFeedGenerator(
		logging.CreateLogger(zap.DebugLevel),
		configureFriends(),
		configureUsers(),
		configureMedia(),
		configureHealth(),
	)
	m.Run()
}

func TestFeedGenerator_GenerateFeedForUser(t *testing.T) {
	posts, err := gen.GenerateFeedForUser(context.Background(), 1, "")

	if err != nil {
		t.Errorf("Failed to generate feed with error: %v", err)
	}

	// first there are three posts on the same day, assert that they are in order of closest friends
	if posts[0].Metadata["userID"] != "5" {
		t.Errorf("Expected post to be by user 5, got: %v", posts[0].Metadata["userID"])
	}
	if posts[1].Metadata["userID"] != "2" {
		t.Errorf("Expected post to be by user 2, got: %v", posts[1].Metadata["userID"])
	}
	if posts[2].Metadata["userID"] != "3" {
		t.Errorf("Expected post to be by user 3, got: %v", posts[2].Metadata["userID"])
	}
	if posts[3].Metadata["userID"] != "5" {
		t.Errorf("Expected post to be by user 5, got: %v", posts[3].Metadata["userID"])
	}
}

func TestFeedGenerator_GenerateFeedForNegativeUser(t *testing.T) {
	posts, err := gen.GenerateFeedForUser(context.Background(), 4, "")

	if err != nil {
		t.Errorf("Failed to generate feed with error: %v", err)
	}

	// this user is quite negative, the first post should be mental health based, and then every 5 after that
	if posts[0].Metadata["userID"] != "-1" {
		t.Errorf("Expected first post to be by user -1, got: %v", posts[0].Metadata["userID"])
	}
	if posts[5].Metadata["userID"] != "-1" {
		t.Errorf("Expected first post to be by user -1, got: %v", posts[5].Metadata["userID"])
	}
	if posts[10].Metadata["userID"] != "-1" {
		t.Errorf("Expected first post to be by user -1, got: %v", posts[10].Metadata["userID"])
	}

	for _, post := range posts {
		fmt.Println(post)
	}
}
