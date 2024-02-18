package imageset

type ImageSet struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Images        []string  `json:"images"`
	AverageColors []float64 `json:"averageColor"`
}
