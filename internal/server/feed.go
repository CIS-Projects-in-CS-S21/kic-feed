package server

import (
	pbfeed "github.com/kic/feed/pkg/proto/feed"
)

type FeedService struct {
	pbfeed.UnimplementedFeedServer
}

func NewFeedService() *FeedService {
	return &FeedService{}
}

func (f *FeedService) GenerateFeedForUser(req *pbfeed.GenerateFeedForUserRequest, stream pbfeed.Feed_GenerateFeedForUserServer) error {
	return nil
}