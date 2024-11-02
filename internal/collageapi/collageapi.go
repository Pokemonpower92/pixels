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
	r.RegisterRoute("POST /images/sets", imageSetHandler.CreateImageSet)
	r.RegisterRoute("GET /images/sets", imageSetHandler.GetImageSets)
	r.RegisterRoute("GET /images/sets/{id}", imageSetHandler.GetImageSetById)

	tiRepo, err := repository.NewTagrgetImageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	targetImageHandler := handler.NewTargetImageHandler(tiRepo)
	r.RegisterRoute("POST /images/targets", targetImageHandler.CreateTargetImage)
	r.RegisterRoute("GET /images/targets", targetImageHandler.GetTargetImages)
	r.RegisterRoute("GET /images/targets/{id}", targetImageHandler.GetTargetImageById)

	acRepo, err := repository.NewAverageColorRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	averageColorHandler := handler.NewAverageColorHandler(acRepo)
	r.RegisterRoute("POST /images/averagecolors", averageColorHandler.CreateAverageColor)
	r.RegisterRoute("GET /images/averagecolors", averageColorHandler.GetAverageColors)
	r.RegisterRoute("GET /images/averagecolors/{id}", averageColorHandler.GetAverageColorById)

	cRepo, err := repository.NewCollageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	collageHandler := handler.NewCollageHandler(cRepo)
	r.RegisterRoute("POST /images/collages", collageHandler.CreateCollage)
	r.RegisterRoute("GET /images/collages", collageHandler.GetCollages)
	r.RegisterRoute("GET /images/collages/{id}", collageHandler.GetCollageById)

	s := server.NewImageSetServer(r)
	s.Start()
}
