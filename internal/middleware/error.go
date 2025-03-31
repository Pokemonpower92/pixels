package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/logger"
	"github.com/pokemonpower92/collagegenerator/internal/response"
)

func Error() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					l, _ := logger.GetRequestLogger(r)
					err := response.WriteResponse(
						w,
						http.StatusInternalServerError,
						errors.New("Unknown server error."),
					)
					if err != nil {
						l.Error(fmt.Sprintf("Error writing response: %+v %+v", err, w))
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
