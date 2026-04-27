package utils

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var LogChannel = make(chan string, 100)

// InitLogger initializes the global logger based on log level
func InitLogger(logLevel string) error {
    // ... existing initialization code ...
    // (Keeping it simple to just add the channel logic)
    return initLoggerInternal(logLevel)
}

func initLoggerInternal(logLevel string) error {
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
		Encoding:         "console",
		EncoderConfig:    zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}

func broadcast(level, msg string) {
    select {
    case LogChannel <- fmt.Sprintf("[%s] %s", level, msg):
    default:
        // Channel full, drop log
    }
}

// Info logs info level message
func Info(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(message, fields...)
        broadcast("INFO", message)
	}
}

// Debug logs debug level message
func Debug(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(message, fields...)
        broadcast("DEBUG", message)
	}
}

// Warn logs warning level message
func Warn(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(message, fields...)
        broadcast("WARN", message)
	}
}

// Error logs error level message
func Error(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(message, fields...)
        broadcast("ERROR", message)
	}
}

// Fatal logs fatal level message and exits
func Fatal(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(message, fields...)
        broadcast("FATAL", message)
	}
}
