package util

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func LogInfo(ctx context.Context, message string) {
	slog.InfoContext(ctx, message, "traceID", ctx.Value("traceID"))
}

func LogError(ctx context.Context, message string, err error) {
	slog.ErrorContext(ctx, message, "traceID", ctx.Value("traceID"), "error", err)
}

func CreateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
