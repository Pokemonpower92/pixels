package store

import (
	_ "image/jpeg"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

// Store is an interface that defines the methods
// for retrieving and storing images.
type Store interface {
	GetFile(id uuid.UUID) (io.Reader, error)
	PutFile(id uuid.UUID, reader io.Reader) error
}

type StoreFunc = func() Store

func NewStore(l *slog.Logger) Store {
	return NewLocalStore(os.Getenv("STORE_DIRECTORY"), l)
}
