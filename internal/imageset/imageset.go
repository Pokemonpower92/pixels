package imageset

import "image/color"

type ImageSet struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	AverageColors []color.RGBA `json:"averageColor"`
}
