package common

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger(env string) {
	level := slog.LevelInfo
	if env == "dev" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	Logger = slog.New(handler)
}

func LogInfo(msg string, args ...any) {
	if Logger != nil {
		Logger.Info(msg, args...)
	}
}

func LogDebug(msg string, args ...any) {
	if Logger != nil {
		Logger.Debug(msg, args...)
	}
}

func LogWarn(msg string, args ...any) {
	if Logger != nil {
		Logger.Warn(msg, args...)
	}
}

func LogError(msg string, args ...any) {
	if Logger != nil {
		Logger.Error(msg, args...)
	}
}
