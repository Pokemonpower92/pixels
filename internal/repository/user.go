// Repositories are thin wrappers over the generated sqlc queries.

package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"
)

// UserModeler is an interface for sqlc queries
type UserModeler interface {
	Get(userName string) (*sqlc.User, error)
	Create(sqlc.CreateUserParams) (*sqlc.User, error)
}

type UserRepository struct {
	client *pgxpool.Pool
	logger *slog.Logger
	ctx    context.Context
	q      *sqlc.Queries
}

func NewUserRepository(
	client *pgxpool.Pool,
	ctx context.Context,
) (*UserRepository, error) {
	q := sqlc.New(client)
	return &UserRepository{
		client: client,
		logger: slog.Default(),
		ctx:    ctx,
		q:      q,
	}, nil
}

func (ir *UserRepository) Close() {
	ir.client.Close()
}

func (ir *UserRepository) Get(userName string) (*sqlc.User, error) {
	User, err := ir.q.GetUser(ir.ctx, userName)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (ir *UserRepository) Create(query sqlc.CreateUserParams) (*sqlc.User, error) {
	User, err := ir.q.CreateUser(ir.ctx, query)
	if err != nil {
		return nil, err
	}
	return User, nil
}
