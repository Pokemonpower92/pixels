package router

import (
	"log/slog"
	"net/http"

	"github.com/pokemonpower92/pixels/internal/logger"
	"github.com/pokemonpower92/pixels/internal/middleware"
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

func (r *Router) RegisterRoute(path string, handler ApiFunc) {
	handlerFunc := makeHttpHandler(handler)
	stdMiddleware := middleware.New(
		middleware.Context(),
		middleware.Timer(),
	)
	handlerFunc = stdMiddleware.Use(handlerFunc)
	r.Mux.HandleFunc(path, handlerFunc.ServeHTTP)
}
