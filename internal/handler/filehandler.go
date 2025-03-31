package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/response"
)

func GetFiles(w http.ResponseWriter, _ *http.Request) error {
	return nil
}

func GetFileById(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Getting File by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	store := datastore.NewStore(l)
	image, err := store.GetFile(id)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, image)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Got File: %s", id))
	return nil
}

func StoreFile(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Storing File")
	id := uuid.New()
	store := datastore.NewStore(l)
	if err := store.PutFile(id, r.Body); err != nil {
		return err
	}
	l.Info("Stored File")
	response.WriteResponse(w, http.StatusCreated, id)
	return nil
}
