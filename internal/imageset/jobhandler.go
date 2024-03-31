package imageset

import (
	"log"

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
	l     *log.Logger
	cache *Cache
	db    *db.ImageSetDB
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
	im, err := jh.cache.GetImageSet(job.ImagesetID)
	if err != nil {
		jh.l.Printf("Failed to get imageset from cache: %s", err)
	}

	if im == nil {
		jh.l.Printf("Imageset not found in cache, generating it")

		g := NewGenerator(job)
		im, err := g.Generate()
		if err != nil {
			jh.l.Printf("Failed to generate imageset: %s", err)
		}

		jh.l.Printf("Generated imageset: %v", im.Name)

		err = jh.cache.SetImageSet(im)
		if err != nil {
			jh.l.Printf("Failed to set imageset in cache: %s", err)
		}

	} else {
		jh.l.Printf("Got imageset from cache: %v", im.Name)
	}
}

// convert Id to int
