package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type ImageSetHandler struct {
	l    *log.Logger
	repo repository.ISRepo
}

func NewImageSetHandler(repo repository.ISRepo) *ImageSetHandler {
	l := log.New(log.Writer(), "", log.LstdFlags)
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

func (ish *ImageSetHandler) CreateImageSet(w http.ResponseWriter, r *http.Request) error {
	ish.l.Printf("Creating ImageSet")
	var req sqlc.CreateImageSetParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	imageSet, err := ish.repo.Create(req)
	if err != nil {
		return err
	}
	ish.l.Printf("Created ImageSet with id: %s", imageSet.ID)
	response.WriteResponse(w, http.StatusCreated, imageSet)
	return nil
}
