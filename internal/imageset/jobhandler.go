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
	jh.l.Printf("Got imageset from cache: %v", im)
}
