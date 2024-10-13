package repository

import (
	"fmt"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/domain"
)

type Repository[O any] interface {
	Get(id int) (*O, bool)
	GetAll() ([]*O, bool)
	Create(obj *O) error
	Update(id int, obj *O) (*O, error)
	Delete(id int) error
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
	ISRepo Repository[domain.ImageSet]
	TIRepo Repository[domain.TargetImage]
)
