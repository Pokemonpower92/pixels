package datastore

import (
	"image"
	_ "image/jpeg"
	"os"

	"github.com/google/uuid"
)

// Store is an interface that defines the methods
// for retrieving and storing images.
type Store interface {
	GetImage(id uuid.UUID) (*image.RGBA, error)
	PutImage(id uuid.UUID, image *image.RGBA) error
}

type StoreFunc = func() Store

func NewStore() Store {
	configMap := make(map[string]StoreFunc)
	configMap["LOCAL"] = func() Store {
		return NewLocalStore(os.Getenv("STORE_ROOT"))
	}
	return configMap[os.Getenv("STORE_HOST")]()
}
