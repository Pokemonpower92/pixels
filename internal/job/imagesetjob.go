package job

type ImageSetJob struct {
	ImagesetID  string `json:"imageset_id"`
	BucketName  string `json:"bucket_name"`
	Description string `json:"description"`
}

func NewImageSetJob(jobJson map[string]interface{}) *ImageSetJob {
	return &ImageSetJob{
		ImagesetID:  jobJson["imageset_id"].(string),
		BucketName:  jobJson["bucket_name"].(string),
		Description: jobJson["description"].(string),
	}
}
