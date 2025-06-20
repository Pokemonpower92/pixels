package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/pokemonpower92/pixels/internal/logger"
	"github.com/pokemonpower92/pixels/internal/response"
	"github.com/pokemonpower92/pixels/internal/session"
)

// Auth validates session and adds user_id to context
func Auth(sessionizer session.Sessionizer) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l, _ := logger.GetRequestLogger(r)
			cookie, err := r.Cookie("session_id")
			if err != nil {
				l.Info("No session cookie found")
				response.WriteErrorResponse(w, 401, err)
				return
			}
			userID, ok := sessionizer.FindSession(cookie.Value)
			if !ok {
				l.Info("Invalid session", "session_id", cookie.Value, "error", err)
				response.WriteErrorResponse(w, 401, errors.New("Invalid Session"))
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", userID.String())
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

// Helper function to get user ID from context
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value("user_id").(string)
	return userID, ok
}
