package imageset

import (
	"log"
)

type Job struct {
	ImagesetID string `json:"imageset_id"`
	BucketName string `json:"bucket_name"`
	Path       string `json:"path"`
}

func NewJob(jobJson map[string]interface{}) *Job {
	return &Job{
		ImagesetID: jobJson["imageset_id"].(string),
		BucketName: jobJson["bucket_name"].(string),
		Path:       jobJson["path"].(string),
	}
}

type JobHandler struct {
	l     *log.Logger
	cache *Cache
}

func NewJobHandler(l *log.Logger) *JobHandler {
	return &JobHandler{
		l:     l,
		cache: NewCache(l),
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

		g := NewGenerator(jh.l, job)
		im, err := g.Generate()
		if err != nil {
			jh.l.Printf("Failed to generate imageset: %s", err)
		}

		jh.l.Printf("Generated imageset: %v", im)

		err = jh.cache.SetImageSet(im)
		if err != nil {
			jh.l.Printf("Failed to set imageset in cache: %s", err)
		}

	} else {
		jh.l.Printf("Got imageset from cache: %v", im)
	}
}
