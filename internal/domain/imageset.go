package domain

import "image/color"

type AverageColor struct {
	FileName string     `json:"file_name"`
	Color    color.RGBA `json:"color"`
}

type ImageSet struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	AverageColors []*AverageColor `json:"average_colors"`
}
