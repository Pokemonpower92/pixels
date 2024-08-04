package collageapi

import (
	"log"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func Start() {
	r := router.NewRouter()
	l := log.New(log.Writer(), "imagesethandler: ", log.LstdFlags)
	c := config.NewPostgresConfig()
	repo, err := repository.NewImageSetRepository(c)
	if err != nil {
		panic(err)
	}
	imageSetHandler := handler.NewImageSetHandler(l, repo)
	imageSetsHandler := handler.NewImageSetsHandler(l, repo)
	r.RegisterHandler("/imagesets", imageSetsHandler)
	r.RegisterHandler("/imagesets/{id}", imageSetHandler)

	s := server.NewImageSetServer(r)
	s.Start()
}
