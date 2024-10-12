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
	config.LoadEnvironmentVariables()
	r := router.NewRouter()
	l := log.New(log.Writer(), "imagesethandler: ", log.LstdFlags)
	c := config.NewPostgresConfig()
	repo, err := repository.NewImageSetRepository(c)
	if err != nil {
		panic(err)
	}
	imageSetHandler := handler.NewImageSetHandler(l, repo)
	r.RegisterHandler("GET /images/sets", imageSetHandler.GetImageSets)
	r.RegisterHandler("GET /images/sets/{id}", imageSetHandler.GetImageSetById)

	s := server.NewImageSetServer(r)
	s.Start()
}
