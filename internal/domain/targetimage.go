package domain

type TargetImage struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Path        string `json:"path"`
	Description string `json:"description"`
}
