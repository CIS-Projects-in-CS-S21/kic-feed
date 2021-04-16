package main

import (
	"context"
	"fmt"
	pbcommon "github.com/kic/feed/pkg/proto/common"
	pbfeed "github.com/kic/feed/pkg/proto/feed"
	"google.golang.org/grpc/metadata"
	"io"
	"log"

	"google.golang.org/grpc"

	pbusers "github.com/kic/feed/pkg/proto/users"
)

func postsShouldBeInOrder(authctx context.Context, uid int64, client pbfeed.FeedClient) []*pbcommon.File {
	expectedOrder := map[int]string{
		1: "Makefile4",
		2: "Makefile5",
		3: "Makefile1",
		4: "Makefile2",
		6: "Makefile6",
		7: "Makefile3",
	}

	feedRes, err := client.GenerateFeedForUser(authctx, &pbfeed.GenerateFeedForUserRequest{UserID: uid})

	if err != nil {
		log.Fatalf("fail to get user ID: %v", err)
	}

	var files []*pbcommon.File

	for {
		recv, err := feedRes.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Got an error gneerating feed that should not occur: %v", err)
		}
		fmt.Printf("feed res: %v\n", recv.FileInfo)
		files = append(files, recv.FileInfo)
	}

	for idx, file := range files {
		if val, ok := expectedOrder[idx]; ok {
			if val != file.FileName {
				log.Fatalf("File out of order at index %v! %v\nExpected %v", idx, file, expectedOrder[idx])
			}
		}

	}

	log.Printf("postsShouldBeInOrder: Success")
	return files
}

func postsShouldHaveMentalHealthInjections(posts []*pbcommon.File) {
	if posts[0].Metadata["userID"] != "150" {
		log.Fatalf("Mental health post expected at index 0 but not seen")
	}

	if posts[5].Metadata["userID"] != "150" {
		log.Fatalf("Mental health post expected at index 5 but not seen")
	}
	log.Printf("postsShouldHaveMentalHealthInjections: Success")
}

func main() {
	conn, err := grpc.Dial("test.api.keeping-it-casual.com:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pbfeed.NewFeedClient(conn)
	userClient := pbusers.NewUsersClient(conn)

	tokRes, err := userClient.GetJWTToken(context.Background(), &pbusers.GetJWTTokenRequest{
		Username: "testuser",
		Password: "testpass",
	})

	if err != nil {
		log.Fatalf("fail to get token: %v", err)
	}

	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := userClient.GetUserByUsername(ctx, &pbusers.GetUserByUsernameRequest{Username: "testuser"})
	if err != nil {
		log.Fatalf("fail to get user ID: %v", err)
	}

	uid := res.User.UserID

	posts := postsShouldBeInOrder(ctx, uid, client)

	postsShouldHaveMentalHealthInjections(posts)

	log.Print("Success")
}
