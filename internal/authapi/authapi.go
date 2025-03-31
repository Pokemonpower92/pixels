package authapi

import (
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func Start() {
	r := router.NewRouter()

	r.RegisterRoute("POST /authenticate", handler.Authenticate)
	r.RegisterRoute("POST /authorize", handler.Authorize)

	s := server.NewAuthServer(r)
	s.Start()
}
