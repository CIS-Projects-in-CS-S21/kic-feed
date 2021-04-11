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
		0:"Makefile4",
		1:"Makefile5",
		2:"Makefile1",
		3:"Makefile2",
		4:"Makefile6",
		5:"Makefile3",
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

		if expectedOrder[idx] != file.FileName {
			log.Fatalf("File out of order! %v\nExpected %v", file, expectedOrder[idx] )
		}
	}

	return files
}

func postsShouldHaveMentalHealthInjections() {

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

	postsShouldBeInOrder(ctx, uid, client)

	postsShouldHaveMentalHealthInjections()

	log.Print("Success")
}
