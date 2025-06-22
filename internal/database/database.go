package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pokemonpower92/pixels/config"
)

// GetConnectionString renders a connection string from the given config
func GetConnectionString(config *config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
}

// NewDatabase returns a new pgxpool.Pool
func NewDatabase(
	pgConfig *config.DBConfig,
	ctx context.Context,
) (*pgxpool.Pool, error) {
	connString := GetConnectionString(pgConfig)
	return pgxpool.New(
		ctx,
		connString,
	)
}
