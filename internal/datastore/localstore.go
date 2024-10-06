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
	if err != nil {
		s.logger.Printf("Failed to open directory: %s", err)
		return nil, err
	}
	defer f.Close()

}
