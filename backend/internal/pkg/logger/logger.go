package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

type contextKey string

const LoggerKey contextKey = "logger"

func New(level, format string, w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}
	var l slog.Level
	switch strings.ToLower(level) {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	var handler slog.Handler
	switch strings.ToLower(format) {
	case "text":
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{Level: l})
	default:
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{Level: l})
	}
	return slog.New(handler)
}

func FromCtx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
