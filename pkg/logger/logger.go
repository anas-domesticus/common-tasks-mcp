package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new zap logger instance
// If verbose is true, uses Development config with Debug level
// If verbose is false, uses Production config with Info level
func New(verbose bool) (*zap.Logger, error) {
	var config zap.Config

	if verbose {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	return config.Build()
}

// NewNop creates a no-op logger that discards all logs
func NewNop() *zap.Logger {
	return zap.NewNop()
}
