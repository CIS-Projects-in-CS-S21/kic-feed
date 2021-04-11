// This provides a script to upload files for the target test users friends so that a feed can be generated

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pbcommon "github.com/kic/feed/pkg/proto/common"
	pbmedia "github.com/kic/feed/pkg/proto/media"
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

	buffer, err := ioutil.ReadFile("Makefile")

	if err != nil {
		log.Fatal("cannot read file: ", err)
	}

	mediaClient := pbmedia.NewMediaStorageClient(conn)

	req := &pbmedia.UploadFileRequest{
		FileInfo: &pbcommon.File{
			FileName:     "Makefile1",
			FileLocation: "test",
			Metadata: map[string]string{
				"userID":  "30",
				"ext":     "txt",
				"caption": "test",
			},
			DateStored: &pbcommon.Date{
				Year:  2021,
				Month: 3,
				Day:   2,
			},
		},
		File: buffer,
	}

	_, err = mediaClient.UploadFile(ctx, req)

	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req = &pbmedia.UploadFileRequest{
		FileInfo: &pbcommon.File{
			FileName:     "Makefile2",
			FileLocation: "test",
			Metadata: map[string]string{
				"userID":  "30",
				"ext":     "txt",
				"caption": "test",
			},
			DateStored: &pbcommon.Date{
				Year:  2021,
				Month: 2,
				Day:   1,
			},
		},
		File: buffer,
	}

	_, err = mediaClient.UploadFile(ctx, req)

	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req = &pbmedia.UploadFileRequest{
		FileInfo: &pbcommon.File{
			FileName:     "Makefile3",
			FileLocation: "test",
			Metadata: map[string]string{
				"userID":  "30",
				"ext":     "txt",
				"caption": "test",
			},
			DateStored: &pbcommon.Date{
				Year:  2020,
				Month: 2,
				Day:   1,
			},
		},
		File: buffer,
	}

	_, err = mediaClient.UploadFile(ctx, req)

	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req = &pbmedia.UploadFileRequest{
		FileInfo: &pbcommon.File{
			FileName:     "Makefile4",
			FileLocation: "test",
			Metadata: map[string]string{
				"userID":  "31",
				"ext":     "txt",
				"caption": "test",
			},
			DateStored: &pbcommon.Date{
				Year:  2021,
				Month: 5,
				Day:   1,
			},
		},
		File: buffer,
	}

	_, err = mediaClient.UploadFile(ctx, req)

	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req = &pbmedia.UploadFileRequest{
		FileInfo: &pbcommon.File{
			FileName:     "Makefile5",
			FileLocation: "test",
			Metadata: map[string]string{
				"userID":  "31",
				"ext":     "txt",
				"caption": "test",
			},
			DateStored: &pbcommon.Date{
				Year:  2021,
				Month: 4,
				Day:   1,
			},
		},
		File: buffer,
	}

	_, err = mediaClient.UploadFile(ctx, req)

	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req = &pbmedia.UploadFileRequest{
		FileInfo: &pbcommon.File{
			FileName:     "Makefile6",
			FileLocation: "test",
			Metadata: map[string]string{
				"userID":  "31",
				"ext":     "txt",
				"caption": "test",
			},
			DateStored: &pbcommon.Date{
				Year:  2020,
				Month: 2,
				Day:   1,
			},
		},
		File: buffer,
	}

	_, err = mediaClient.UploadFile(ctx, req)

	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}
}
