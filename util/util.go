package util

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

type TodoItem struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

type TodoItemId struct {
	ID int `json:"id"`
}

type TodoRequest struct {
	Action string   `json:"action"`
	Item   TodoItem `json:"item"`
	ID     int      `json:"id"`
	Resp   chan TodoItem
}

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
