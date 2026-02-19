package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global logger instance.
var Log *zap.Logger

// Init initializes the global Zap logger based on environment and log level.
func Init(env, level string) error {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set log level
	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = logger
	zap.ReplaceGlobals(logger)
	return nil
}

// Sync flushes any buffered log entries. Call on shutdown.
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// Fatal logs a message and exits the process.
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
	os.Exit(1)
}
