package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

type loggerKey string

const requestLoggerKey loggerKey = "logger"

type RequestLogger slog.Logger

func StoreRequestLogger(ctx context.Context, rl *slog.Logger) context.Context {
	return context.WithValue(ctx, requestLoggerKey, rl)
}

func GetRequestLogger(r *http.Request) (*slog.Logger, bool) {
	logger, ok := r.Context().Value(requestLoggerKey).(*slog.Logger)
	return logger, ok
}

func NewRequestLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey && len(groups) == 0 {
					return slog.Attr{}
				}
				return a
			},
		}),
	)
}
