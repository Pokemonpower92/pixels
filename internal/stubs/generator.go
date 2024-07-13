package stubs

import (
	"github.com/pokemonpower92/collagegenerator/internal/domain"
	"github.com/pokemonpower92/collagegenerator/internal/job"
)

type GeneratorStub struct {
	GenerateFunc func(job *job.ImageSetJob) (*domain.ImageSet, error)
}

func (g *GeneratorStub) Generate(job *job.ImageSetJob) (*domain.ImageSet, error) {
	return g.GenerateFunc(job)
}
