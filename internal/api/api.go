package api

import (
	"context"
	"log/slog"

	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/handler"
	"github.com/pokemonpower92/pixels/internal/repository"
	"github.com/pokemonpower92/pixels/internal/router"
	"github.com/pokemonpower92/pixels/internal/server"
	"github.com/pokemonpower92/pixels/internal/session"
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
	sessionizer := session.NewSessionStore()
	r.RegisterProtectedRoute("GET /images/{id}", sessionizer, h.GetImage)
	r.RegisterProtectedRoute("GET /images", sessionizer, h.GetImages)
	r.RegisterProtectedRoute("POST /images", sessionizer, h.CreateImage)

	userRepo, err := repository.NewUserRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer userRepo.Close()
	authHandler := handler.NewAuthHandler(
		userRepo,
		sessionizer,
		slog.Default(),
	)
	r.RegisterUnprotectedRoute("POST /users", authHandler.CreateUser)
	r.RegisterUnprotectedRoute("POST /login", authHandler.Login)

	r.RegisterUnprotectedRoute("GET /healthcheck", handler.HealthCheck)

	serverConfig := config.NewServerConfig()
	s := server.NewServer(r, serverConfig)
	s.Start()
}
