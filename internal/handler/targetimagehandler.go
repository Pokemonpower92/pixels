package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/repository"
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

func (tih *TargetImageHandler) GetTargetImages(w http.ResponseWriter, _ *http.Request) {
	tih.l.Printf("Getting TargetImages")
	w.Write([]byte("Got TargetImages\n"))
}

func (tih *TargetImageHandler) GetTargetImageById(w http.ResponseWriter, r *http.Request) {
	tih.l.Printf("Getting TargetImage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		tih.l.Printf("Invalid id: %s", err)
		http.Error(w, "Invalid ID:", http.StatusInternalServerError)
		return
	}
	response := fmt.Sprintf("Got TargetImage with id: %d\n", id)
	w.Write([]byte(response))
}
