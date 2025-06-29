package handler

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pokemonpower92/pixels/internal/database"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupTestContainer(ctx context.Context) (*postgres.PostgresContainer, string, error) {
	dbName := "pixels"
	dbUser := "user"
	dbPassword := "password"
	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts("../sqlc/migrations/postgres"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, "", err
	}
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return container, "", err
	}
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		return container, "", err
	}
	if err := database.RunMigration(config); err != nil {
		return container, "", err
	}
	return container, connStr, nil
}

type PostTestCase struct {
	name           string
	shouldError    bool
	expectedError  error
	expectedStatus int
	data           string
}

func PostApiTestFunc(
	t *testing.T,
	endpoint string,
	testCases []PostTestCase,
) func(*testing.T) {
	return func(t *testing.T) {
		for _, test := range testCases {
			resp, err := http.Post(
				endpoint,
				"application/json",
				strings.NewReader(test.data),
			)
			if err != nil {
				t.Errorf("Test %s failed: %s", test.name, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != test.expectedStatus {
				t.Errorf(
					"Test %s failed: expected = %v, got %v",
					test.name,
					test.expectedStatus,
					resp.Status,
				)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if test.shouldError && !strings.Contains(string(body), `"error"`) {
				t.Errorf(
					"Test %s failed: expected error, but got none",
					test.name,
				)
			}
		}
	}
}
