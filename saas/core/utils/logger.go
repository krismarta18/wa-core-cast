package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the global logger based on log level
func InitLogger(logLevel string) error {
	var atomicLevel zap.AtomicLevel

	switch logLevel {
	case "debug":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warning":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		atomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	config := zap.Config{
		Level:            atomicLevel,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	defer Logger.Sync()
	return nil
}

// Info logs info level message
func Info(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(message, fields...)
	}
}

// Debug logs debug level message
func Debug(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(message, fields...)
	}
}

// Warn logs warning level message
func Warn(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(message, fields...)
	}
}

// Error logs error level message
func Error(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(message, fields...)
	}
}

// Fatal logs fatal level message and exits
func Fatal(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(message, fields...)
	}
}
