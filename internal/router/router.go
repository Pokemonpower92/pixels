package router

import (
	"net/http"
)

type Router struct {
	Mux *http.ServeMux
}

func NewRouter() *Router {
	sm := http.NewServeMux()
	return &Router{Mux: sm}
}

func (r *Router) RegisterRoute(path string, handler http.HandlerFunc) {
	r.Mux.HandleFunc(path, handler)
}
