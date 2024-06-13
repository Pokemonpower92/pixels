package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pokemonpower92/collagecommon/types"
)

type ImageSetRepository struct {
	client *pgxpool.Pool
	logger *log.Logger
}

func getConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("IMAGESET_DB"),
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

func (isr *ImageSetRepository) Get(id int) (*types.ImageSet, bool) {
	return &types.ImageSet{
		Name: "test",
	}, true
}

func (isr *ImageSetRepository) Create(is *types.ImageSet) error {
	return nil
}

func (isr *ImageSetRepository) Update(id int, is *types.ImageSet) (*types.ImageSet, error) {
	return &types.ImageSet{
		Name: "test",
	}, nil
}

func (isr *ImageSetRepository) Delete(id int) (*types.ImageSet, error) {
	return &types.ImageSet{
		Name: "test",
	}, nil
}
