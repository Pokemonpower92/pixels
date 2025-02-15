package main

import (
	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/service"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

func main() {
	config.LoadEnvironmentVariables()
	collageId := uuid.MustParse("3245a68b-4e7b-4b77-bf79-bbeafb93e413")
	collageImage := sqlc.CollageImage{
		ID:        uuid.New(),
		CollageID: collageId,
	}
	service.GenerateCollage(&collageImage)
}
