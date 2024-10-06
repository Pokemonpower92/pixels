package datastore

import (
	"image"
	"log"
	"os"
)

type LocalStore struct {
	Directory string
	logger    *log.Logger
}

func NewLocalStore(directory string) *LocalStore {
	return &LocalStore{
		Directory: directory,
		logger:    log.New(log.Writer(), "localstore ", log.LstdFlags),
	}
}

func (s *LocalStore) GetImages() ([]*image.RGBA, error) {
	s.logger.Printf("Reading images from directory: %s", s.Directory)
	f, err := os.Open(s.Directory)
}
