package jobhandler

import (
	"log"
	"strconv"

	"github.com/pokemonpower92/imagesetservice/internal/generator"
	"github.com/pokemonpower92/imagesetservice/internal/job"
	"github.com/pokemonpower92/imagesetservice/internal/repository"
)

// JobHandler handles the processing of image set jobs.
type JobHandler struct {
	logger     *log.Logger
	repository repository.ISRepo
}

func NewJobHandler(
	repository repository.ISRepo,
	log *log.Logger) *JobHandler {
	return &JobHandler{
		logger:     log,
		repository: repository,
	}
}

func (jobHandler *JobHandler) GenerateImageSet(job *job.Job) error {
	generatorLogger := log.New(log.Writer(), "generator: ", log.Flags())
	generator := generator.NewImageSetGenerator(job, generatorLogger)

	imageSet, err := generator.Generate(job)
	if err != nil {
		jobHandler.logger.Printf("Failed to generate imageset: %s", err)
		return err
	}

	jobHandler.logger.Printf("Generated imageset: %v", imageSet.Name)
	err = jobHandler.repository.Create(imageSet)
	if err != nil {
		jobHandler.logger.Printf("Failed to add new imageset to repository: %s", err)
		return err
	} else {
		jobHandler.logger.Printf("Added imageset to repository: %v", imageSet.Name)
	}

	return nil
}

func (jobHandler *JobHandler) HandleJob(job *job.Job) error {
	jobHandler.logger.Printf("Handling job: %v", job)

	idAsInteger, err := strconv.Atoi(job.ImagesetID)
	if err != nil {
		jobHandler.logger.Printf("Failed to convert imageset id to int: %s", err)
		return err
	}

	imageSet, ok := jobHandler.repository.Get(idAsInteger)
	if !ok {
		jobHandler.logger.Printf("Imageset not found in repository: %v", job.ImagesetID)
		err := jobHandler.GenerateImageSet(job)
		if err != nil {
			jobHandler.logger.Printf("Failed to generate imageset: %s", err)
			return err
		}
	} else {
		jobHandler.logger.Printf("Imageset found in repository: %v", imageSet.Name)
	}
	return nil
}
