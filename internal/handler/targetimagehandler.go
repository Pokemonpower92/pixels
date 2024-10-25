package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/utils"
)

type TargetImageHandler struct {
	l    *log.Logger
	repo repository.TIRepo
}

func NewTargetImageHandler(repo repository.TIRepo) *TargetImageHandler {
	l := log.New(log.Writer(), "TargetImageHandler: ", log.LstdFlags)
	return &TargetImageHandler{
		l:    l,
		repo: repo,
	}
}

func (tih *TargetImageHandler) GetTargetImages(w http.ResponseWriter, _ *http.Request) error {
	tih.l.Printf("Getting TargetImages")
	utils.WriteJson(w, http.StatusOK, "Got all target images")
	return nil
}

func (tih *TargetImageHandler) GetTargetImageById(w http.ResponseWriter, r *http.Request) error {
	tih.l.Printf("Getting TargetImage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		tih.l.Printf("Invalid id: %s", err)
		return err
	}
	tih.l.Printf("Got target image for id: %s", id)
	utils.WriteJson(w, http.StatusOK, id)
	return nil
}
