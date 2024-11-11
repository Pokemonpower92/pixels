package handler

import "github.com/pokemonpower92/collagegenerator/internal/router"

type Handler interface {
	RegisterRoutes(r *router.Router)
}
