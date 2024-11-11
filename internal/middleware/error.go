package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/response"
)

func Error() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("[%s] [%s] [%s]\n", "ERROR", r.URL, err)
					err := response.WriteResponse(
						w,
						http.StatusInternalServerError,
						errors.New("Unknown server error."),
					)
					if err != nil {
						log.Printf("Error writing response: %+v %+v", err, w)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
