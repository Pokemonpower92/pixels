package datastore

import (
	"image"
	_ "image/jpeg"
)

// Store is an interface that defines the methods for retrieving image sets.
type Store interface {
	GetImages() ([]*image.RGBA, error)
}
