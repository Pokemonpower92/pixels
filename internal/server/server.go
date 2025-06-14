package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/router"
)

type Server struct {
	server *http.Server
	router *router.Router
}

func (s *Server) Start() {
	log.Printf("Starting server on %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func NewServer(router *router.Router, config *config.ServerConfig) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler: router.Mux,
	}

	return &Server{
		server: server,
		router: router,
	}
}
