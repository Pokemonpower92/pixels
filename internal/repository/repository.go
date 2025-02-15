package repository

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type Repository[O, R any] interface {
	Get(id uuid.UUID) (*O, error)
	GetAll() ([]*O, error)
	Create(req R) (*O, error)
	Update(id uuid.UUID, req R) (*O, error)
	Delete(id uuid.UUID) error
}

func GetConnectionString(config *config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
}

type (
	ISRepo Repository[sqlc.ImageSet, sqlc.CreateImageSetParams]
	TIRepo Repository[sqlc.TargetImage, sqlc.CreateTargetImageParams]
	ACRepo Repository[sqlc.AverageColor, sqlc.CreateAverageColorParams]
	CRepo  Repository[sqlc.Collage, sqlc.CreateCollageParams]
	CIRepo Repository[sqlc.CollageImage, uuid.UUID]
)
