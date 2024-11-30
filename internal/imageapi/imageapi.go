package imageapi

import (
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func Start() {
	config.LoadEnvironmentVariables()

	r := router.NewRouter()
	h := handler.NewImageHandler()
	r.RegisterRoute("POST /images/{id}", h.StoreImage)

	s := server.NewImageServer(r)
	s.Start()
}
