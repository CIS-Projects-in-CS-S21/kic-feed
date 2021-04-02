package server

import (
	"context"
	pbfeed "github.com/kic/feed/pkg/proto/feed"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type FeedService struct {
	pbfeed.UnimplementedFeedServer

	feedGen *FeedGenerator
	logger  *zap.SugaredLogger
}

func NewFeedService(feedGen *FeedGenerator, logger *zap.SugaredLogger) *FeedService {
	return &FeedService{
		feedGen: feedGen,
		logger:  logger,
	}
}

func (f *FeedService) GenerateFeedForUser(
	req *pbfeed.GenerateFeedForUserRequest,
	stream pbfeed.Feed_GenerateFeedForUserServer,
) error {
	ctx := stream.Context()
	headers, success := metadata.FromIncomingContext(ctx)

	if !success {
		return status.Errorf(codes.Unauthenticated, "Send token along with request")
	}

	header := headers["authorization"][0]

	posts, err := f.feedGen.GenerateFeedForUser(context.Background(), req.UserID, header)

	if err != nil {
		return status.Errorf(codes.Internal, "Failed to generate feed: %v", err)
	}

	for _, post := range posts {
		err = stream.Send(
			&pbfeed.GenerateFeedForUserResponse{FileInfo: post},
		)

		if err != nil {
			return status.Errorf(codes.Unavailable, "Failed to send post stream")
		}
	}

	return nil
}
