package middleware

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/internal/logger"
)

// Inject required objects into the request's context.
func Context() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			baseLogger := logger.NewRequestLogger()
			requestLogger := baseLogger.WithGroup("request").With(
				slog.String("id", uuid.NewString()),
				slog.String("method", r.Method),
				slog.String("path", r.RequestURI),
			)
			ctx := logger.StoreRequestLogger(r.Context(), requestLogger)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		},
		)
	}
}
