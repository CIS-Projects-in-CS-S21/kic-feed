package logging

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CreateLogger - Create a new logger instance
func CreateLogger(level zapcore.Level) *zap.SugaredLogger {
	var config zap.Config
	// Setup Logging
	if os.Getenv("PRODUCTION") != "" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = zap.NewAtomicLevelAt(level)
	loggerMgr, err := config.Build()
	if err != nil {
		log.Fatalf("Couldn't start zap logger: %v", err)
	}

	logger := loggerMgr.Sugar()

	return logger
}
