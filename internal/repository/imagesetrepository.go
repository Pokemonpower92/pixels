package repository

import (
	"context"
	"fmt"
	"image/color"
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
	connString := getConnectionString(postgresConfig)
	client, err := pgxpool.New(
		context.Background(),
		connString,
	)
	if err != nil {
		return nil, err
	}
	logger := log.New(
		log.Writer(),
		"ImageSetRepository: ",
		log.LstdFlags,
	)
	return &ImageSetRepository{client: client, logger: logger}, nil
}

func (isr *ImageSetRepository) Close() {
	isr.client.Close()
}

func (isr *ImageSetRepository) Get(id int) (*domain.ImageSet, bool) {
	var imageSet domain.ImageSet
	var averageColors []color.RGBA

	err := isr.client.QueryRow(
		context.Background(),
		"SELECT name, description FROM imagesets WHERE id = $1",
		id,
	).Scan(&imageSet.Name, &imageSet.Description)
	if err != nil {
		isr.logger.Printf("Failed to get imageset: %s", err)
		return nil, false
	}

	rows, err := isr.client.Query(
		context.Background(),
		"SELECT red, green, blue, alpha FROM average_colors WHERE imageset_id = $1",
		id,
	)
	if err != nil {
		isr.logger.Printf("Failed to get average colors: %s", err)
		return nil, false
	}
	defer rows.Close()

	for rows.Next() {
		var color color.RGBA
		err := rows.Scan(&color.R, &color.G, &color.B, &color.A)
		if err != nil {
			isr.logger.Printf("Failed to scan average color: %s", err)
			return nil, false
		}
		averageColors = append(averageColors, color)
	}

	if err := rows.Err(); err != nil {
		isr.logger.Printf("Failed to iterate over average colors: %s", err)
		return nil, false
	}

	imageSet.AverageColors = averageColors

	return &imageSet, true
}

func (isr *ImageSetRepository) Create(is *domain.ImageSet) error {
	var id int
	err := isr.client.QueryRow(
		context.Background(),
		"INSERT INTO imagesets (name, description) VALUES ($1, $2) RETURNING id",
		is.Name,
		is.Description,
	).Scan(&id)
	if err != nil {
		isr.logger.Printf("Failed to create imageset: %s", err)
		return err
	}

	for _, color := range is.AverageColors {
		_, err := isr.client.Exec(
			context.Background(),
			"INSERT INTO average_colors (imageset_id, red, green, blue, alpha) VALUES ($1, $2, $3, $4, $5)",
			id,
			color.R,
			color.G,
			color.B,
			color.A,
		)
		if err != nil {
			isr.logger.Printf("Failed to create average color: %s", err)
			return err
		}
	}

	return nil
}

func (isr *ImageSetRepository) Update(id int, is *domain.ImageSet) (*domain.ImageSet, error) {
	tx, err := isr.client.Begin(context.Background())
	if err != nil {
		isr.logger.Printf("Failed to begin transaction: %s", err)
		return nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		"UPDATE imagesets SET name = $1, description = $2 WHERE id = $3",
		is.Name,
		is.Description,
		id,
	)
	if err != nil {
		isr.logger.Printf("Failed to update imageset: %s", err)
		tx.Rollback(context.Background())
		return nil, err
	}

	_, err = tx.Exec(
		context.Background(),
		"DELETE FROM average_colors WHERE imageset_id = $1",
		id,
	)
	if err != nil {
		isr.logger.Printf("Failed to delete average colors: %s", err)
		tx.Rollback(context.Background())
		return nil, err
	}

	for _, color := range is.AverageColors {
		_, err := tx.Exec(
			context.Background(),
			"INSERT INTO average_colors (imageset_id, red, green, blue, alpha) VALUES ($1, $2, $3, $4, $5)",
			id,
			color.R,
			color.G,
			color.B,
			color.A,
		)
		if err != nil {
			isr.logger.Printf("Failed to create average color: %s", err)
			tx.Rollback(context.Background())
			return nil, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		isr.logger.Printf("Failed to commit transaction: %s", err)
		tx.Rollback(context.Background())
		return nil, err
	}

	return is, nil
}

func (isr *ImageSetRepository) Delete(id int) error {
	tx, err := isr.client.Begin(context.Background())
	if err != nil {
		isr.logger.Printf("Failed to begin transaction: %s", err)
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		"DELETE FROM average_colors WHERE imageset_id = $1",
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
