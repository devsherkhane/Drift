package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global logger instance
var Log *zap.Logger

// SetupLogging initializes the global zap logger
func SetupLogging() {
	var config zap.Config

	// Determine if running in production or development
	env := os.Getenv("ENV")
	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human readable timestamps
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Colorize levels in dev
	}

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	Log = logger
}

// Sync flushes any buffered log entries. Should be called before application exit.
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}