package repository

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type Repository[O any] interface {
	Get(id uuid.UUID) (*O, error)
	GetAll() ([]*O, error)
	Create(obj *O) error
	Update(id uuid.UUID, obj *O) (*O, error)
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
	ISRepo Repository[sqlc.Imageset]
	TIRepo Repository[sqlc.TargetImage]
)
