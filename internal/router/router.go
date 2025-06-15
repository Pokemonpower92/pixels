package router

import (
	"log/slog"
	"net/http"

	"github.com/pokemonpower92/pixels/internal/logger"
	"github.com/pokemonpower92/pixels/internal/middleware"
	"github.com/pokemonpower92/pixels/internal/session"
)

type ApiFunc func(http.ResponseWriter, *http.Request, *slog.Logger)

type ApiError struct {
	Error string
}

func makeHttpHandler(h ApiFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, _ := logger.GetRequestLogger(r)
		h(w, r, l)
	})
}

type Router struct {
	Mux *http.ServeMux
}

func NewRouter() *Router {
	sm := http.NewServeMux()
	return &Router{Mux: sm}
}

func (r *Router) RegisterProtectedRoute(
	path string,
	sessionizer session.Sessionizer,
	handler ApiFunc,
) {
	handlerFunc := makeHttpHandler(handler)
	stdMiddleware := middleware.New(
		middleware.Context(),
		middleware.Auth(sessionizer),
		middleware.Timer(),
	)
	handlerFunc = stdMiddleware.Use(handlerFunc)
	r.Mux.HandleFunc(path, handlerFunc.ServeHTTP)
}

func (r *Router) RegisterUnprotectedRoute(
	path string,
	handler ApiFunc,
) {
	handlerFunc := makeHttpHandler(handler)
	stdMiddleware := middleware.New(
		middleware.Context(),
		middleware.Timer(),
	)
	handlerFunc = stdMiddleware.Use(handlerFunc)
	r.Mux.HandleFunc(path, handlerFunc.ServeHTTP)
}
