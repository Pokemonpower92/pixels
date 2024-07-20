package main

import (
	"log"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func main() {
	l := log.New(log.Writer(), "imagesethandler: ", log.LstdFlags)
	c := config.NewPostgresConfig()
	repo, err := repository.NewImageSetRepository(c)
	if err != nil {
		panic(err)
	}
	h := handler.NewImageSetHandler(l, repo)
	r := router.NewRouter(h)
	s := server.NewImageSetServer(r)
	s.Start()
}
