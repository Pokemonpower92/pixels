package jobhandler

import (
	"log"
	"strconv"

	"github.com/pokemonpower92/collagegenerator/internal/generator"
	"github.com/pokemonpower92/collagegenerator/internal/job"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
)

// ISJobHandler handles the processing of image set jobs.
type ISJobHandler struct {
	logger     *log.Logger
	repository repository.ISRepo
}

func NewISJobHandler(
	repository repository.ISRepo,
	log *log.Logger) *ISJobHandler {
	return &ISJobHandler{
		logger:     log,
		repository: repository,
	}
}

func (isjh *ISJobHandler) generateImageSet(job *job.Job) error {
	generatorLogger := log.New(log.Writer(), "generator: ", log.Flags())
	generator := generator.NewImageSetGenerator(job, generatorLogger)

	imageSet, err := generator.Generate(job)
	if err != nil {
		isjh.logger.Printf("Failed to generate imageset: %s", err)
		return err
	}

	isjh.logger.Printf("Generated imageset: %v", imageSet.Name)
	err = isjh.repository.Create(imageSet)
	if err != nil {
		isjh.logger.Printf("Failed to add new imageset to repository: %s", err)
		return err
	} else {
		isjh.logger.Printf("Added imageset to repository: %v", imageSet.Name)
	}

	return nil
}

func (isjh *ISJobHandler) HandleJob(job *job.Job) error {
	isjh.logger.Printf("Handling job: %v", job)

	idAsInteger, err := strconv.Atoi(job.ImagesetID)
	if err != nil {
		isjh.logger.Printf("Failed to convert imageset id to int: %s", err)
		return err
	}

	imageSet, ok := isjh.repository.Get(idAsInteger)
	if !ok {
		isjh.logger.Printf("Imageset not found in repository: %v", job.ImagesetID)
		err := isjh.generateImageSet(job)
		if err != nil {
			isjh.logger.Printf("Failed to generate imageset: %s", err)
			return err
		}
	} else {
		isjh.logger.Printf("Imageset found in repository: %v", imageSet.Name)
	}
	return nil
}
