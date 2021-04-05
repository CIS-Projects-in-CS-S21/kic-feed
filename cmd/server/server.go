package main

import (
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/kic/feed/internal/server"
	"github.com/kic/feed/pkg/logging"
	pbfeed "github.com/kic/feed/pkg/proto/feed"
)

func main() {
	IsProduction := os.Getenv("PRODUCTION") != ""
	var logger *zap.SugaredLogger
	var connectionURL string
	if IsProduction {
		logger = logging.CreateLogger(zapcore.InfoLevel)
		connectionURL = "keeping-it-casual.com:50051"
	} else {
		logger = logging.CreateLogger(zapcore.DebugLevel)
		connectionURL = "test.keeping-it-casual.com:50051"
	}

	ListenAddress := ":" + os.Getenv("PORT")

	listener, err := net.Listen("tcp", ListenAddress)

	if err != nil {
		logger.Fatalf("Unable to listen on %v: %v", ListenAddress, err)
	}

	grpcServer := grpc.NewServer()

	if err != nil {
		logger.Fatalf("Unable connect to db %v", err)
	}

	conn, err := grpc.Dial(connectionURL, grpc.WithInsecure())

	if err != nil {
		logger.Fatalf("fail to dial: %v", err)
	}

	feedGen := server.NewFeedGenerator(
		logger,
		server.NewFriendClientWrapper(conn),
		server.NewUserClientWrapper(conn),
		server.NewMediaClientWrapper(conn),
		server.NewHealthClientWrapper(conn),
	)

	serv := server.NewFeedService(feedGen, logger)

	pbfeed.RegisterFeedServer(grpcServer, serv)

	go func() {
		defer listener.Close()
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	defer grpcServer.Stop()

	// the server is listening in a goroutine so hang until we get an interrupt signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
}
