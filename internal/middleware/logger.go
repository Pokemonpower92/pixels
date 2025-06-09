package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pokemonpower92/pixels/internal/logger"
)

func Logger() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(w, r)
			elapsedTime := time.Since(startTime)
			logger, _ := logger.GetRequestLogger(r)
			logger.Info(fmt.Sprintf("request complete in %s", elapsedTime))
		},
		)
	}
}
