package middleware

import (
	"fmt"
	"net/http"

	"github.com/pokemonpower92/pixels/internal/logger"
)

func Error() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					l, _ := logger.GetRequestLogger(r)
					l.Error(fmt.Sprintf("Error serving request: %s", err))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
