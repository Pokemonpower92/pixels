package api

import (
	"context"
	"log/slog"

	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/handler"
	"github.com/pokemonpower92/pixels/internal/repository"
	"github.com/pokemonpower92/pixels/internal/router"
	"github.com/pokemonpower92/pixels/internal/server"
)

func Start() {
	r := router.NewRouter()
	c := config.NewPostgresConfig()
	ctx := context.Background()
	imageRepo, err := repository.NewImageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer imageRepo.Close()

	h := handler.NewImageHandler(
		imageRepo,
		*config.NewResolutionConfig(),
		slog.Default(),
	)
	r.RegisterRoute("GET /images/{id}", h.GetImage)
	r.RegisterRoute("GET /images", h.GetImages)
	r.RegisterRoute("POST /images", h.CreateImage)

	r.RegisterRoute("GET /healthcheck", handler.HealthCheck)

	serverConfig := config.NewServerConfig()
	s := server.NewAuthServer(r, serverConfig)
	s.Start()
}
