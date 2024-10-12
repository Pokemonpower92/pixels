package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/pokemonpower92/collagegenerator/internal/domain"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
)

type isRepository = repository.Repository[domain.ImageSet]

type ImageSetHandler struct {
	l    *log.Logger
	repo isRepository
}

func NewImageSetHandler(l *log.Logger, repo isRepository) *ImageSetHandler {
	return &ImageSetHandler{l: l, repo: repo}
}

func (ish *ImageSetHandler) GetImageSets(w http.ResponseWriter, _ *http.Request) {
	ish.l.Printf("Getting ImageSets")
	imageSets, ok := ish.repo.GetAll()
	if !ok {
		ish.l.Printf("Failed to get ImageSets")
		http.Error(w, "Error getting ImageSets", http.StatusInternalServerError)
		return
	}
	ish.l.Printf("Found %d ImageSets.", len(imageSets))

	encoder := json.NewEncoder(w)
	err := encoder.Encode(imageSets)
	if err != nil {
		ish.l.Printf("Failed to encode ImageSets: %s", err)
		http.Error(w, "Error encoding ImageSets", http.StatusInternalServerError)
	}
}

func (ish *ImageSetHandler) GetImageSetById(w http.ResponseWriter, r *http.Request) {
	ish.l.Printf("Getting ImageSet by ID")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ish.l.Printf("Invalid ID: %s", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	imageSet, ok := ish.repo.Get(id)
	if !ok {
		ish.l.Printf("ImageSet not found")
		http.Error(w, "ImageSet not found", http.StatusNotFound)
		return
	}
	ish.l.Printf("Found ImageSet: %v", imageSet)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(imageSet)
	if err != nil {
		ish.l.Printf("Failed to encode ImageSet: %s", err)
		http.Error(w, "Error encoding ImageSet", http.StatusInternalServerError)
	}
}
