package router

import (
	"log"
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/middleware"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandler(h ApiFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			panic(err)
		}
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
		middleware.Logger(log.New(log.Writer(), "", log.LstdFlags)),
		middleware.Error(),
	)
	handlerFunc = stdMiddleware.Use(handlerFunc)
	r.Mux.HandleFunc(path, handlerFunc.ServeHTTP)
}
