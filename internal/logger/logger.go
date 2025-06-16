package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

const (
	ContextKeyLogger = "logger"
)

func New(ctx context.Context, output io.Writer, level slog.Level) (context.Context, *slog.Logger) {

	handler := &ContextHandler{
		Handler: slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: level,
		}),
	}

	logger := slog.New(handler)
	return context.WithValue(ctx, ContextKeyLogger, logger), logger
}

func FromContext(ctx context.Context) *slog.Logger {
	logger := ctx.Value(ContextKeyLogger)

	if logger == nil {
		_, logger = New(ctx, os.Stdout, slog.LevelDebug)
	}
	return logger.(*slog.Logger)
}

type ContextHandler struct {
	slog.Handler
}
