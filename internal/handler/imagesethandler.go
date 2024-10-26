package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
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
	response.WriteResponse(w, http.StatusOK, imageSets)
	return nil
}

func (ish *ImageSetHandler) GetImageSetById(w http.ResponseWriter, r *http.Request) error {
	ish.l.Printf("Getting ImageSet by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	imageSet, err := ish.repo.Get(id)
	if err != nil {
		return err
	}
	ish.l.Printf("Found ImageSet: %v", imageSet)
	response.WriteResponse(w, http.StatusOK, imageSet)
	return nil
}
