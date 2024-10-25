package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/utils"
)

type ImageSetHandler struct {
	l    *log.Logger
	repo repository.ISRepo
}

func NewImageSetHandler(repo repository.ISRepo) *ImageSetHandler {
	l := log.New(log.Writer(), "ImageSetHandler: ", log.LstdFlags)
	return &ImageSetHandler{l: l, repo: repo}
}

func (ish *ImageSetHandler) GetImageSets(w http.ResponseWriter, _ *http.Request) error {
	ish.l.Printf("Getting ImageSets")
	imageSets, err := ish.repo.GetAll()
	if err != nil {
		return err
	}
	ish.l.Printf("Found %d ImageSets", len(imageSets))
	encoder := json.NewEncoder(w)
	err = encoder.Encode(imageSets)
	if err != nil {
		ish.l.Printf("Failed to encode ImageSets: %s", err)
		return err
	}
	utils.WriteJson(w, http.StatusOK, imageSets)
	return nil
}

func (ish *ImageSetHandler) GetImageSetById(w http.ResponseWriter, r *http.Request) error {
	ish.l.Printf("Getting ImageSet by ID")
	id := uuid.MustParse(r.PathValue("id"))
	imageSet, err := ish.repo.Get(id)
	if err != nil {
		return err
	}
	ish.l.Printf("Found ImageSet: %v", imageSet)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(imageSet)
	if err != nil {
		return err
	}
	utils.WriteJson(w, http.StatusOK, imageSet)
	return nil
}
