package router

import (
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/utils"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, h *http.Request) {
		if err := f(w, h); err != nil {
			utils.WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
		}
	}
}

type Router struct {
	Mux *http.ServeMux
}

func NewRouter() *Router {
	sm := http.NewServeMux()
	return &Router{Mux: sm}
}

func (r *Router) RegisterRoute(path string, handler apiFunc) {
	r.Mux.HandleFunc(path, makeHttpHandler(handler))
}
