package router

import (
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/handler"
)

type Router struct {
	Mux *http.ServeMux
}

func NewRouter(imageSetHandler handler.Handler) *Router {
	sm := http.NewServeMux()
	sm.Handle("/imagesets/{id}", imageSetHandler)

	return &Router{Mux: sm}
}
