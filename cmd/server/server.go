package main

import (
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/kic/feed/pkg/logging"
	"github.com/kic/feed/internal/server"
	pbfeed "github.com/kic/feed/pkg/proto/feed"
)

func main() {
	IsProduction := os.Getenv("PRODUCTION") != ""
	var logger *zap.SugaredLogger
	if IsProduction {
		logger = logging.CreateLogger(zapcore.InfoLevel)
	} else {
		logger = logging.CreateLogger(zapcore.DebugLevel)
	}

	ListenAddress := ":" + os.Getenv("PORT")

	listener, err := net.Listen("tcp", ListenAddress)
	if err != nil {
		logger.Fatalf("Unable to listen on %v: %v", ListenAddress, err)
	}

	grpcServer := grpc.NewServer()

	if err != nil {
		logger.Fatalf("Unable connect to db %v",  err)
	}


	if err != nil {
		logger.Fatalf("Unable migrate tables to db %v",  err)
	}

	serv := server.NewFeedService()

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