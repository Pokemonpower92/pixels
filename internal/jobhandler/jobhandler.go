package jobhandler

import (
	"log"
	"strconv"

	"github.com/pokemonpower92/collagecommon/types"
	"github.com/pokemonpower92/imagesetservice/internal/generator"
	"github.com/pokemonpower92/imagesetservice/internal/job"
	"github.com/pokemonpower92/imagesetservice/internal/repository"
)

// JobHandler handles the processing of image set jobs.
type JobHandler struct {
	logger     *log.Logger
	repository repository.Repository[types.ImageSet]
	generator  generator.Generator
}

// NewJobHandler creates a new JobHandler instance.
func NewJobHandler() *JobHandler {
	log := log.New(log.Writer(), "jobhandler ", log.LstdFlags)
	repository, err := repository.NewImageSetRepository()
	if err != nil {
		log.Fatalf("Failed to create imageset repository: %s", err)
	}
	return &JobHandler{
		logger:     log,
		repository: repository,
	}
}

// HandleJob handles the processing of a single job.
func (jobHandler *JobHandler) HandleJob(job *job.Job) error {
	jobHandler.logger.Printf("Handling job: %v", job)
	idAsInteger, err := strconv.Atoi(job.ImagesetID)
	if err != nil {
		jobHandler.logger.Printf("Failed to convert imageset id to int: %s", err)
		return err
	}
	_, ok := jobHandler.repository.Get(idAsInteger)
	if !ok {
		jobHandler.generator = generator.NewImageSetGenerator(job)
		imageSet, err := jobHandler.generator.Generate(job)
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
	}
	return nil
}
