package collageapi

import (
	"context"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func Start() {
	config.LoadEnvironmentVariables()

	r := router.NewRouter()
	c := config.NewPostgresConfig()
	ctx := context.Background()
	isRepo, err := repository.NewImageSetRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	imageSetHandler := handler.NewImageSetHandler(isRepo)
	r.RegisterRoute("GET /images/sets", imageSetHandler.GetImageSets)
	r.RegisterRoute("GET /images/sets/{id}", imageSetHandler.GetImageSetById)

	tiRepo, err := repository.NewTagrgetImageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	targetImageHandler := handler.NewTargetImageHandler(tiRepo)
	r.RegisterRoute("GET /images/targets", targetImageHandler.GetTargetImages)
	r.RegisterRoute("GET /images/targets/{id}", targetImageHandler.GetTargetImageById)

	s := server.NewImageSetServer(r)
	s.Start()
}
