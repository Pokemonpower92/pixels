package api

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/auth"
	"github.com/pokemonpower92/pixels/internal/handler"
	"github.com/pokemonpower92/pixels/internal/repository"
	"github.com/pokemonpower92/pixels/internal/router"
	"github.com/pokemonpower92/pixels/internal/server"
	"github.com/pokemonpower92/pixels/internal/session"
)

// Start is the entrypoint for the api.
func Start() {
	ctx := context.Background()
	connString := config.ConnString()
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	imageRepo, err := repository.NewImageRepository(db, ctx)
	if err != nil {
		panic(err)
	}
	defer imageRepo.Close()
	h := handler.NewImageHandler(
		imageRepo,
		*config.NewResolutionConfig(),
		slog.Default(),
	)

	privateKey, err := auth.GetPrivateKey(config.PrivateKeyPem())
	if err != nil {
		panic(err)
	}
	jwtManager := auth.JwtManager{PrivateKey: privateKey}
	sessionizer := session.NewJWTSessionizer(&jwtManager)

	r := router.NewRouter()
	r.RegisterProtectedRoute("GET /images/{id}", sessionizer, h.GetImage)
	r.RegisterProtectedRoute("GET /images", sessionizer, h.GetImages)
	r.RegisterProtectedRoute("POST /images", sessionizer, h.CreateImage)

	userRepo, err := repository.NewUserRepository(db, ctx)
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
