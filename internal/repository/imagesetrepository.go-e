package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pokemonpower92/imagesetservice/config"
	"github.com/pokemonpower92/imagesetservice/internal/domain"
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
