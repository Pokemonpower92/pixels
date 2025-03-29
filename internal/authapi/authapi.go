package authapi

import (
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func Start() {
	h := handler.NewAuthHandler()
	r := router.NewRouter()

	r.RegisterRoute("POST /authenticate", h.Authenticate)
	r.RegisterRoute("POST /authorize", h.Authorize)

	s := server.NewAuthServer(r)
	s.Start()
}
