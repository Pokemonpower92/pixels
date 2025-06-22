//go:build integration

package handler

// import (
// 	"context"
// 	"log/slog"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/pokemonpower92/pixels/config"
// 	"github.com/pokemonpower92/pixels/internal/auth"
// 	"github.com/pokemonpower92/pixels/internal/repository"
// 	"github.com/pokemonpower92/pixels/internal/router"
// 	"github.com/pokemonpower92/pixels/internal/session"
// )

// func TestMain(m *testing.M) {
// 	code := m.Run()
// 	os.Exit(code)
// }

// func setupTestServer(t *testing.T, db string) *httptest.Server {
// 	r := router.NewRouter()
// 	c := config.NewPostgresConfig()
// 	ctx := context.Background()
// 	privateKey, err := auth.GetPrivateKey(config.PrivateKeyPem())
// 	if err != nil {
// 		panic(err)
// 	}
// 	jwtManager := auth.JwtManager{PrivateKey: privateKey}
// 	sessionizer := session.NewJWTSessionizer(&jwtManager)
// 	userRepo, err := repository.NewUserRepository(c, ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	authHandler := NewAuthHandler(
// 		userRepo,
// 		sessionizer,
// 		slog.Default(),
// 	)
// 	r.RegisterUnprotectedRoute("POST /users", authHandler.CreateUser)
// 	r.RegisterUnprotectedRoute("POST /login", authHandler.Login)
// 	return httptest.NewServer(r.Mux)
// }

// func TestAuthIntegration(t *testing.T) {
// 	//db := setupTestDB(t)
// 	//defer cleanupTestDB(t, db)
// 	//server := setupTestServer(t, db)

// 	t.Run("create user flow", func(t *testing.T) {
// 		// Test user creation
// 	})

// 	t.Run("login flow", func(t *testing.T) {
// 		// Test login with created user
// 	})
// }
