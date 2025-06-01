package main

import (
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func start() {
	r := router.NewRouter()

	r.RegisterRoute("POST /files", handler.StoreFile)
	r.RegisterRoute("GET /files/{id}", handler.GetFileById)

	serverConfig := config.NewServerConfig()
	s := server.NewAuthServer(r, serverConfig)
	s.Start()
}

func main() {
	start()
}
