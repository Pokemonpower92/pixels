package jobhandler

import "github.com/pokemonpower92/collagegenerator/internal/job"

type JobHandler interface {
	HandleJob(*job.Job) error
}
