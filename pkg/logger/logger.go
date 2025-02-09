package logger

import (
	"io"
	"log"
	"log/slog"
)

func MustNew(lvl string, w io.Writer) *slog.Logger {
	switch lvl {
	case "debug":
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "info":
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "warn":
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelWarn}))
	case "error":
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelError}))
	default:
		log.Fatalf("unknown log level: %s", lvl)
		return nil
	}
}
