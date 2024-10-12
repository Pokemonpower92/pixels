package job

type ImageSetJob struct {
	ImagesetID  string `json:"imageset_id"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

func NewImageSetJob(jobJson map[string]interface{}) *ImageSetJob {
	return &ImageSetJob{
		ImagesetID:  jobJson["imageset_id"].(string),
		Path:        jobJson["path"].(string),
		Description: jobJson["description"].(string),
	}
}
