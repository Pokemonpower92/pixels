package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pokemonpower92/collagegenerator/internal/domain"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
)

type tiRepository = repository.Repository[domain.ImageSet]

type TargetImageHandler struct {
	l    *log.Logger
	repo tiRepository
}

func NewTargetImageHandler(l *log.Logger, repo tiRepository) *TargetImageHandler {
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		tih.l.Printf("Invalid id: %s", err)
		http.Error(w, "Invalid ID:", http.StatusInternalServerError)
		return
	}
	response := fmt.Sprintf("Got TargetImage with id: %d\n", id)
	w.Write([]byte(response))
}
