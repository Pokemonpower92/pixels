package job

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
