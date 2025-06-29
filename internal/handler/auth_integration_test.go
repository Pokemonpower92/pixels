//go:build integration

package handler

import (
	"context"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/auth"
	"github.com/pokemonpower92/pixels/internal/database"
	"github.com/pokemonpower92/pixels/internal/repository"
	"github.com/pokemonpower92/pixels/internal/router"
	"github.com/pokemonpower92/pixels/internal/session"
	"github.com/testcontainers/testcontainers-go"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setupTestServer(ctx context.Context, db *pgxpool.Pool) *httptest.Server {
	r := router.NewRouter()
	privateKey, err := auth.GetPrivateKey(config.PrivateKeyPem())
	if err != nil {
		panic(err)
	}
	jwtManager := auth.JwtManager{PrivateKey: privateKey}
	sessionizer := session.NewJWTSessionizer(&jwtManager)
	userRepo, err := repository.NewUserRepository(db, ctx)
	if err != nil {
		panic(err)
	}
	authHandler := NewAuthHandler(
		userRepo,
		sessionizer,
		slog.Default(),
	)
	r.RegisterUnprotectedRoute("POST /users", authHandler.CreateUser)
	r.RegisterUnprotectedRoute("POST /login", authHandler.Login)
	return httptest.NewServer(r.Mux)
}

func TestAuthIntegration(t *testing.T) {
	ctx := context.Background()
	container, connStr, err := SetupTestContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}()
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.RunMigration(config); err != nil {
		t.Fatal(err)
	}

	db, err := pgxpool.New(ctx, connStr)
	server := setupTestServer(ctx, db)
	defer server.Close()
	baseUrl := server.URL

	testCases := []PostTestCase{
		{
			name:           "Success: create new user",
			shouldError:    false,
			expectedError:  nil,
			expectedStatus: 200,
			data:           `{"user_name": "testUser", "password": "testPass"}`,
		},
		{
			name:           "Fail: duplicate user",
			shouldError:    true,
			expectedError:  nil,
			expectedStatus: 500,
			data:           `{"user_name": "testUser", "password": "testPass"}`,
		},
	}

	t.Run("Create user flow", PostApiTestFunc(t, baseUrl+"/users", testCases))
}
