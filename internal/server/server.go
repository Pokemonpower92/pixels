package server

import (
	"log"
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/router"
)

type ImageSetServer struct {
	server *http.Server
	router *router.Router
}

func NewImageSetServer(router *router.Router) *ImageSetServer {
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router.Mux,
	}

	return &ImageSetServer{
		server: server,
		router: router,
	}
}

func (iss *ImageSetServer) Start() {
	log.Printf("Starting server on %s", iss.server.Addr)
	err := iss.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
