package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(logger *log.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			logger.Printf("[%s] [%s]\n", r.Method, r.URL)
			h.ServeHTTP(w, r)
			elapsedTime := time.Since(startTime)
			logger.Printf(
				"[%s] [%s] [%s]\n",
				r.Method,
				r.URL,
				elapsedTime,
			)
		},
		)
	}
}
