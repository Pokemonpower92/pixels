package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pokemonpower92/imagesetservice/internal/domain"
)

type ImageSetRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
}

func getConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

func NewImageSetRepository() (*ImageSetRepository, error) {
	connString := getConnectionString()
	client, err := pgxpool.Connect(
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

func (isr *ImageSetRepository) Get(id int) (*domain.ImageSet, bool) {
	return &domain.ImageSet{
		Name: "test",
	}, true
}

func (isr *ImageSetRepository) Create(is *domain.ImageSet) error {
	return nil
}

func (isr *ImageSetRepository) Update(id int, is *domain.ImageSet) (*domain.ImageSet, error) {
	return &domain.ImageSet{
		Name: "test",
	}, nil
}

func (isr *ImageSetRepository) Delete(id int) (*domain.ImageSet, error) {
	return &domain.ImageSet{
		Name: "test",
	}, nil
}
