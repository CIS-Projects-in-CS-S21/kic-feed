// This provides a script to upload files for the target test users friends so that a feed can be generated

package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pbhealth "github.com/kic/feed/pkg/proto/health"
	pbusers "github.com/kic/feed/pkg/proto/users"
)

//func

func main() {
	conn, err := grpc.Dial("test.api.keeping-it-casual.com:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	//client := pbfeed.NewFeedClient(conn)
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

	healthClient := pbhealth.NewHealthTrackingClient(conn)

	resp, err := healthClient.AddHealthDataForUser(ctx, &pbhealth.AddHealthDataForUserRequest{
		UserID: uid,
		NewEntry: &pbhealth.MentalHealthLog{
			LogDate:     nil,
			Score:       -5,
			JournalName: "asdfasdf",
			UserID:      uid,
		},
	})
	if err != nil {
		return
	}
	fmt.Printf("%v\n", resp.Success)
}
