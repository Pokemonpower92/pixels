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

type ImageSetRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
}

func NewImageSetRepository(postgresConfig *config.DBConfig) (*ImageSetRepository, error) {
	logger := log.New(
		log.Writer(),
		"ImageSetRepository: ",
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
	return &ImageSetRepository{client: client, logger: logger}, nil
}

func (isr *ImageSetRepository) Close() {
	isr.client.Close()
}

func (isr *ImageSetRepository) Get(id uuid.UUID) (*sqlc.Imageset, bool) {
	ctx := context.Background()
	q := sqlc.New(isr.client)
	imageSet, err := q.GetImageset(ctx, id)
	if err != nil {
		isr.logger.Printf("Failed to get all imagesets: %s", err)
		return nil, false
	}
	return imageSet, true
}

func (isr *ImageSetRepository) GetAll() ([]*sqlc.Imageset, bool) {
	ctx := context.Background()
	q := sqlc.New(isr.client)
	imageSets, err := q.ListImagesets(ctx)
	if err != nil {
		isr.logger.Printf("Failed to get all imagesets: %s", err)
		return nil, false
	}
	return imageSets, true
}

func (isr *ImageSetRepository) Create(is *sqlc.Imageset) error {
	isr.logger.Printf("Create not implemented for imageset.\n")
	return errors.New("Not implemented")
}

func (isr *ImageSetRepository) Update(id uuid.UUID, is *sqlc.Imageset) (*sqlc.Imageset, error) {
	isr.logger.Printf("Update not implemented for imageset.\n")
	return nil, errors.New("Not implemented")
}

func (isr *ImageSetRepository) Delete(id uuid.UUID) error {
	isr.logger.Printf("Delete not implemented for imageset.\n")
	return errors.New("Not implemented")
}
