package job

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
