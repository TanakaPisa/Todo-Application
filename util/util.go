package util

import (
	"context"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func LogInfo(ctx context.Context, message string) {
	Logger.InfoContext(ctx, message, "traceID", ctx.Value("traceID"))
}

func LogError(ctx context.Context, message string, err error) {
	Logger.ErrorContext(ctx, message, "traceID", ctx.Value("traceID"), "error", err)
}
