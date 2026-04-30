package utils

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var LogChannel = make(chan string, 100)

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
		Encoding:         "console",
		EncoderConfig:    zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
		},
		OutputPaths:      []string{"stdout", GetDataPath("debug_log.txt")},
		ErrorOutputPaths: []string{"stderr", GetDataPath("debug_log.txt")},
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
	}
}

func Info(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(message, fields...)
		broadcast("INFO", message)
	}
}

func Debug(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(message, fields...)
		broadcast("DEBUG", message)
	}
}

func Warn(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(message, fields...)
		broadcast("WARN", message)
	}
}

func Error(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(message, fields...)
		broadcast("ERROR", message)
	}
}

func Fatal(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(message, fields...)
		broadcast("FATAL", message)
	}
}
