package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
	ctx    context.Context
	q      *sqlc.Queries
}

func NewCollageRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*CollageRepository, error) {
	logger := log.New(
		log.Writer(),
		"CollageRepository: ",
		log.LstdFlags,
	)
	connString := GetConnectionString(postgresConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &CollageRepository{
		client: client,
		logger: logger,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (cr *CollageRepository) Close() {
	cr.client.Close()
}

func (cr *CollageRepository) Get(id uuid.UUID) (*sqlc.Collage, error) {
	collage, err := cr.q.GetCollage(cr.ctx, id)
	if err != nil {
		return nil, err
	}
	return collage, nil
}

func (cr *CollageRepository) GetAll() ([]*sqlc.Collage, error) {
	collages, err := cr.q.ListCollages(cr.ctx)
	if err != nil {
		return nil, err
	}
	return collages, nil
}

func (cr *CollageRepository) Create(req sqlc.CreateCollageParams) (*sqlc.Collage, error) {
	collage, err := cr.q.CreateCollage(cr.ctx, req)
	if err != nil {
		return nil, err
	}
	return collage, nil
}

func (cr *CollageRepository) Update(
	id uuid.UUID,
	req sqlc.CreateCollageParams,
) (*sqlc.Collage, error) {
	cr.logger.Printf("Update not implemented for collage")
	return nil, errors.New("Not implemented")
}

func (cr *CollageRepository) Delete(id uuid.UUID) error {
	cr.logger.Printf("Delete not implemented for collage")
	return errors.New("Not implemented")
}
