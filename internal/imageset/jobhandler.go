package imageset

import (
	"log"
	"strconv"

	"github.com/pokemonpower92/collagecommon/db"
	"github.com/pokemonpower92/imagesetservice/config"
)

type Job struct {
	ImagesetID  string `json:"imageset_id"`
	BucketName  string `json:"bucket_name"`
	Description string `json:"description"`
}

func NewJob(jobJson map[string]interface{}) *Job {
	return &Job{
		ImagesetID:  jobJson["imageset_id"].(string),
		BucketName:  jobJson["bucket_name"].(string),
		Description: jobJson["description"].(string),
	}
}

type JobHandler struct {
	l         *log.Logger
	cache     iCache
	db        iDB
	generator iGenerator
}

func NewJobHandler() *JobHandler {
	log := log.New(log.Writer(), "jobhandler ", log.LstdFlags)

	db, err := db.NewImageSetDB(config.NewISDBConfig())
	if err != nil {
		log.Fatalf("Failed to create ImageSetDB: %s", err)
	}

	cache := NewCache()

	return &JobHandler{
		l:     log,
		cache: cache,
		db:    db,
	}
}

func (jh *JobHandler) HandleJob(job *Job) {
	jh.l.Printf("Handling job: %v", job)

	intId, err := strconv.Atoi(job.ImagesetID)
	if err != nil {
		jh.l.Printf("Failed to convert imageset id to int: %s", err)
	}

	is, err := jh.db.GetImageSet(intId)
	if err != nil {
		jh.l.Printf("Failed to get imageset from database: %s", err)
	}

	if is == nil {
		is, err := jh.generator.Generate(job)
		if err != nil {
			jh.l.Printf("Failed to generate imageset: %s", err)
		}

		jh.l.Printf("Generated imageset: %v", is.Name)

		err = jh.db.CreateImageSet(is)
		if err != nil {
			jh.l.Printf("Failed to add new imageset to db: %s", err)
		} else {
			jh.l.Printf("Added imageset to db: %v", is.Name)
		}

		err = jh.db.SetAverageColors(intId, is.AverageColors)
		if err != nil {
			jh.l.Printf("Failed to set average colors: %s", err)
		} else {
			jh.l.Printf("Set average colors for imageset: %v", is.Name)
		}

	} else {
		jh.l.Printf("Got imageset from DB: %v", is.Name)
	}
}
