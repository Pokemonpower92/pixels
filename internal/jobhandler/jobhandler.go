package jobhandler

import (
	"log"
	"strconv"

	"github.com/pokemonpower92/collagecommon/db"
	"github.com/pokemonpower92/imagesetservice/config"
)

// Job represents a job for processing an image set.
type Job struct {
	ImagesetID  string `json:"imageset_id"`
	BucketName  string `json:"bucket_name"`
	Description string `json:"description"`
}

// NewJob creates a new Job instance based on the provided JSON data.
func NewJob(jobJson map[string]interface{}) *Job {
	return &Job{
		ImagesetID:  jobJson["imageset_id"].(string),
		BucketName:  jobJson["bucket_name"].(string),
		Description: jobJson["description"].(string),
	}
}

// JobHandler handles the processing of image set jobs.
type JobHandler struct {
	logger    *log.Logger
	cache     Cache
	db        DB
	generator Generator
}

// NewJobHandler creates a new JobHandler instance.
func NewJobHandler() *JobHandler {
	log := log.New(log.Writer(), "jobhandler ", log.LstdFlags)

	db, err := db.NewImageSetDB(config.NewISDBConfig())
	if err != nil {
		log.Fatalf("Failed to create ImageSetDB: %s", err)
	}

	cache := NewImageSetCache()

	return &JobHandler{
		logger: log,
		cache:  cache,
		db:     db,
	}
}

// HandleJob handles the processing of a single job.
func (jobHandler *JobHandler) HandleJob(job *Job) error {
	jobHandler.logger.Printf("Handling job: %v", job)

	idAsInteger, err := strconv.Atoi(job.ImagesetID)
	if err != nil {
		jobHandler.logger.Printf("Failed to convert imageset id to int: %s", err)
		return err
	}

	imageSet, err := jobHandler.db.GetImageSet(idAsInteger)
	if err != nil {
		jobHandler.logger.Printf("Failed to get imageset from database: %s", err)
		return err
	}

	if imageSet == nil {
		jobHandler.generator = NewImageSetGenerator(job)
		imageSet, err := jobHandler.generator.Generate(job)
		if err != nil {
			jobHandler.logger.Printf("Failed to generate imageset: %s", err)
			return err
		}

		jobHandler.logger.Printf("Generated imageset: %v", imageSet.Name)

		err = jobHandler.db.CreateImageSet(imageSet)
		if err != nil {
			jobHandler.logger.Printf("Failed to add new imageset to db: %s", err)
			return err
		} else {
			jobHandler.logger.Printf("Added imageset to db: %v", imageSet.Name)
		}

		err = jobHandler.db.SetAverageColors(idAsInteger, imageSet.AverageColors)
		if err != nil {
			jobHandler.logger.Printf("Failed to set average colors: %s", err)
			return err
		} else {
			jobHandler.logger.Printf("Set average colors for imageset: %v", imageSet.Name)
		}
	} else {
		jobHandler.logger.Printf("Got imageset from DB: %v", imageSet.Name)
	}

	return nil
}
