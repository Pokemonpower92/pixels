package repository

import (
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ImageSetRepository struct {
    client *pgxpool.Pool
    logger *log.Logger
}

func NewImageSetRepository() (*ImageSetRepository, error) {
    return &ImageSetRepository{}, nil
}

func (isr *ImageSetRepository) Get(id int) (int, error) {
    return id, nil    
}

func (isr *ImageSetRepository) Create(o int) error {
    return nil    
}

func (isr *ImageSetRepository) Update(id int, obj int) (int, error) {
    return id, nil    
}

func (isr *ImageSetRepository) Delete(id int) (int, error) {
    return id, nil    
}

