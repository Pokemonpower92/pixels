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

type TargeImageRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
	ctx    context.Context
	q      *sqlc.Queries
}

func NewTagrgetImageRepository(
	pgConfig *config.DBConfig,
	ctx context.Context,
) (*TargeImageRepository, error) {
	logger := log.New(
		log.Writer(),
		"TargeImageRepository: ",
		log.LstdFlags,
	)
	connString := GetConnectionString(pgConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	q := sqlc.New(client)
	return &TargeImageRepository{
		client: client,
		logger: logger,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (tir *TargeImageRepository) Get(id uuid.UUID) (*sqlc.TargetImage, error) {
	tir.logger.Printf("Get not implemented")
	targetImage, err := tir.q.GetTargetImage(tir.ctx, id)
	if err != nil {
		return nil, err
	}
	return targetImage, nil
}

func (tir *TargeImageRepository) GetAll() ([]*sqlc.TargetImage, error) {
	targetImages, err := tir.q.ListTargetImages(tir.ctx)
	if err != nil {
		return nil, err
	}
	return targetImages, nil
}

func (tir *TargeImageRepository) Create(
	req sqlc.CreateTargetImageParams,
) (*sqlc.TargetImage, error) {
	targetImage, err := tir.q.CreateTargetImage(tir.ctx, req)
	if err != nil {
		return nil, err
	}
	return targetImage, nil
}

func (tir *TargeImageRepository) Update(
	id uuid.UUID,
	req sqlc.CreateTargetImageParams,
) (*sqlc.TargetImage, error) {
	tir.logger.Printf("Update not implemented")
	return nil, errors.New("Update not implemented for TargetImages")
}

func (tir *TargeImageRepository) Delete(id uuid.UUID) error {
	tir.logger.Printf("Delete not implemented")
	return errors.New("Delete not implemented for TargetImages")
}
