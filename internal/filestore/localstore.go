package filestore

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalStore struct {
	Directory string
	logger    *slog.Logger
}

func NewLocalStore(directory string, logger *slog.Logger) *LocalStore {
	return &LocalStore{
		Directory: directory,
		logger:    logger,
	}
}

func (s *LocalStore) GetRGBA(id uuid.UUID) (*image.RGBA, error) {
	f, err := s.GetFile(id)
	if err != nil {
		return nil, err
	}
	im, _, err := image.Decode(f)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to decode image: %s", err))
		return nil, err
	}
	b := im.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Bounds(), im, b.Min, draw.Src)
	return rgba, nil
}

func (s *LocalStore) GetFile(id uuid.UUID) (io.Reader, error) {
	path := fmt.Sprintf("%s/%s", s.Directory, id.String())
	s.logger.Info(fmt.Sprintf("Getting File: %s", path))
	f, err := os.Open(path)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to open File: %s", id.String()))
		return nil, err
	}
	return f, nil
}

func (s *LocalStore) PutFile(id uuid.UUID, reader io.Reader) error {
	dst, err := os.Create(filepath.Join(s.Directory, id.String()))
	if err != nil {
		return err
	}
	defer dst.Close()
	s.logger.Info(fmt.Sprintf("Created File destination: %s", dst.Name()))
	if _, err := io.Copy(dst, reader); err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("Successfully stored File"))
	return nil
}
