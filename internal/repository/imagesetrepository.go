package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/domain"
)

type ImageSetRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
}

func getConnectionString(config *config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
}

func NewImageSetRepository(postgresConfig *config.DBConfig) (*ImageSetRepository, error) {
	logger := log.New(
		log.Writer(),
		"ImageSetRepository: ",
		log.LstdFlags,
	)
	connString := getConnectionString(postgresConfig)
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

func (isr *ImageSetRepository) Get(id int) (*domain.ImageSet, bool) {
	var imageSet domain.ImageSet
	err := isr.client.QueryRow(
		context.Background(),
		"SELECT name, description FROM imagesets WHERE id = $1;",
		id,
	).Scan(&imageSet.Name, &imageSet.Description)
	if err != nil {
		isr.logger.Printf("Failed to get imageset: %s", err)
		return nil, false
	}
	return &imageSet, true
}

func (isr *ImageSetRepository) GetAll() ([]*domain.ImageSet, bool) {
	var imageSets []*domain.ImageSet
	rows, err := isr.client.Query(
		context.Background(),
		`SELECT i.id, i.name, i.description, a.red, a.green, a.blue, a.alpha 
		FROM imagesets i;`,
	)
	if err != nil {
		isr.logger.Printf("Failed to get all imagesets: %s", err)
		return nil, false
	}
	defer rows.Close()
	for rows.Next() {
		var imageSet domain.ImageSet
		if err != nil {
			isr.logger.Printf("Failed to scan imageset: %s", err)
			return nil, false
		}
		imageSets = append(imageSets, &imageSet)
	}
	if err := rows.Err(); err != nil {
		isr.logger.Printf("Failed to iterate over imagesets: %s", err)
		return nil, false
	}
	return imageSets, true
}

func (isr *ImageSetRepository) Create(is *domain.ImageSet) error {
	isr.logger.Printf("Create not implemented for imageset.\n")
	return errors.New("Not implemented")
}

func (isr *ImageSetRepository) Update(id int, is *domain.ImageSet) (*domain.ImageSet, error) {
	isr.logger.Printf("Update not implemented for imageset.\n")
	return nil, errors.New("Not implemented")
}

func (isr *ImageSetRepository) Delete(id int) error {
	tx, err := isr.client.Begin(context.Background())
	if err != nil {
		isr.logger.Printf("Failed to begin transaction: %s", err)
		return err
	}
	_, err = tx.Exec(
		context.Background(),
		"DELETE FROM average_colors WHERE imageset_id = $1;",
		id,
	)
	if err != nil {
		isr.logger.Printf("Failed to delete average colors: %s", err)
		tx.Rollback(context.Background())
		return err
	}
	_, err = tx.Exec(
		context.Background(),
		"DELETE FROM imagesets WHERE id = $1",
		id,
	)
	if err != nil {
		isr.logger.Printf("Failed to delete imageset: %s", err)
		tx.Rollback(context.Background())
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		isr.logger.Printf("Failed to commit transaction: %s", err)
		tx.Rollback(context.Background())
		return err
	}
	return nil
}
