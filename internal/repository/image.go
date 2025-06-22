// Repositories are thin wrappers over the generated sqlc queries.

package repository

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"
)

// ImageModeler is an interface for sqlc queries
type ImageModeler interface {
	Get(sqlc.GetImageParams) (*sqlc.Image, error)
	GetAll(userId uuid.UUID) ([]*sqlc.Image, error)
	Create(sqlc.CreateImageParams) (*sqlc.Image, error)
}

type ImageRepository struct {
	client *pgxpool.Pool
	logger *slog.Logger
	ctx    context.Context
	q      *sqlc.Queries
}

func NewImageRepository(
	client *pgxpool.Pool,
	ctx context.Context,
) (*ImageRepository, error) {
	q := sqlc.New(client)
	return &ImageRepository{
		client: client,
		logger: slog.Default(),
		ctx:    ctx,
		q:      q,
	}, nil
}

func (ir *ImageRepository) Close() {
	ir.client.Close()
}

func (ir *ImageRepository) Get(query sqlc.GetImageParams) (*sqlc.Image, error) {
	Image, err := ir.q.GetImage(ir.ctx, query)
	if err != nil {
		return nil, err
	}
	return Image, nil
}

func (ir *ImageRepository) GetAll(userId uuid.UUID) ([]*sqlc.Image, error) {
	Images, err := ir.q.ListImages(ir.ctx, userId)
	if err != nil {
		return nil, err
	}
	return Images, nil
}

func (ir *ImageRepository) Create(query sqlc.CreateImageParams) (*sqlc.Image, error) {
	Image, err := ir.q.CreateImage(ir.ctx, query)
	if err != nil {
		return nil, err
	}
	return Image, nil
}
