package service

import (
	"fmt"

	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

func CreateCollage(collage *sqlc.Collage) {
	fmt.Printf("Creating collage: %+v\n", collage)
}
