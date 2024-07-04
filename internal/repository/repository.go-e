package repository

import (
	"github.com/pokemonpower92/imagesetservice/internal/domain"
)

type Repository[O any] interface {
	Get(id int) (*O, bool)
	Create(obj *O) error
	Update(id int, obj *O) (*O, error)
	Delete(id int) (*O, error)
}

type ISRepo Repository[domain.ImageSet]
