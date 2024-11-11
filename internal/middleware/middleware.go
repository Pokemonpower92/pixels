package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Wrapper struct {
	wrappers []Middleware
}

func New(middleware ...Middleware) Wrapper {
	m := make([]Middleware, 0)
	for i := range middleware {
		m = append(m, middleware[i])
	}
	return Wrapper{wrappers: m}
}

// Use wraps the given http.Handler in the list of middleware.
// Middleware is called in the order specicfied.
func (w Wrapper) Use(h http.Handler) http.Handler {
	for i := len(w.wrappers) - 1; i >= 0; i-- {
		h = w.wrappers[i](h)
	}
	return h
}
