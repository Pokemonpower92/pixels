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
	ctx    context.Context
	q      *sqlc.Queries
}

func NewImageSetRepository(
	postgresConfig *config.DBConfig,
	ctx context.Context,
) (*ImageSetRepository, error) {
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
	q := sqlc.New(client)
	return &ImageSetRepository{
		client: client,
		logger: logger,
		ctx:    ctx,
		q:      q,
	}, nil
}

func (isr *ImageSetRepository) Close() {
	isr.client.Close()
}

func (isr *ImageSetRepository) Get(id uuid.UUID) (*sqlc.ImageSet, error) {
	imageSet, err := isr.q.GetImageSet(isr.ctx, id)
	if err != nil {
		return nil, err
	}
	return imageSet, nil
}

func (isr *ImageSetRepository) GetAll() ([]*sqlc.ImageSet, error) {
	imageSets, err := isr.q.ListImageSets(isr.ctx)
	if err != nil {
		return nil, err
	}
	return imageSets, nil
}

func (isr *ImageSetRepository) Create(req sqlc.CreateImageSetParams) (*sqlc.ImageSet, error) {
	imageset, err := isr.q.CreateImageSet(isr.ctx, req)
	if err != nil {
		return nil, err
	}
	return imageset, nil
}

func (isr *ImageSetRepository) Update(
	id uuid.UUID,
	req sqlc.CreateImageSetParams,
) (*sqlc.ImageSet, error) {
	isr.logger.Printf("Update not implemented for image set")
	return nil, errors.New("Not implemented")
}

func (isr *ImageSetRepository) Delete(id uuid.UUID) error {
	isr.logger.Printf("Delete not implemented for image set")
	return errors.New("Not implemented")
}
