package main

import (
	"log/slog"
	"os"
)

func logLevel(cfg config) slog.Level {
	switch cfg.logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func newLogger(cfg config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: logLevel(cfg),
	}

	switch cfg.logFormat {
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
}
