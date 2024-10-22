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
}

func NewTagrgetImageRepository(pgConfig *config.DBConfig) (*TargeImageRepository, error) {
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
	return &TargeImageRepository{client: client, logger: logger}, nil
}

func (tir *TargeImageRepository) Get(id uuid.UUID) (*sqlc.TargetImage, bool) {
	tir.logger.Printf("Get not implemented")
	return nil, false
}

func (tir *TargeImageRepository) GetAll() ([]*sqlc.TargetImage, bool) {
	tir.logger.Printf("GetAll not implemented")
	return nil, false
}

func (tir *TargeImageRepository) Create(targetImage *sqlc.TargetImage) error {
	tir.logger.Printf("Create not implemented")
	return errors.New("Create not implemented for TargetImages")
}

func (tir *TargeImageRepository) Update(
	id uuid.UUID,
	targetImage *sqlc.TargetImage,
) (*sqlc.TargetImage, error) {
	tir.logger.Printf("Update not implemented")
	return nil, errors.New("Update not implemented for TargetImages")
}

func (tir *TargeImageRepository) Delete(id uuid.UUID) error {
	tir.logger.Printf("Delete not implemented")
	return errors.New("Delete not implemented for TargetImages")
}
