package client

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// FileRequester is an interface for getting a file from filestore.
type FileRequester interface {
	GetFile(uuid.UUID) (*io.Reader, error)
}

type FileReader struct {
	baseURL string
	logger  *slog.Logger
}

func NewFileReader(baseURL string, logger *slog.Logger) *FileReader {
	return &FileReader{
		baseURL: baseURL,
		logger:  logger,
	}
}

func (fr *FileReader) GetFile(id uuid.UUID) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", fr.baseURL, id))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// FileWriter is an interface for storing a file with the given uuid.
type FileWriter interface {
	PutFile(uuid.UUID, io.Reader) error
}
