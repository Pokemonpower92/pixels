package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CreateTargetImageRequest struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	TargetImageId uuid.UUID `json:"targetimage_id"`
}

type TargetImageHandler struct {
	l     *log.Logger
	repo  repository.TIRepo
	store datastore.Store
}

func NewTargetImageHandler(repo repository.TIRepo) *TargetImageHandler {
	l := log.New(log.Writer(), "", log.LstdFlags)
	store := datastore.NewStore()
	return &TargetImageHandler{
		l:     l,
		repo:  repo,
		store: store,
	}
}

func (tih *TargetImageHandler) GetTargetImages(w http.ResponseWriter, _ *http.Request) error {
	tih.l.Printf("Getting TargetImages")
	targetImages, err := tih.repo.GetAll()
	if err != nil {
		return err
	}
	tih.l.Printf("Found %d target images.", len(targetImages))
	response.WriteResponse(w, http.StatusOK, targetImages)
	return nil
}

func (tih *TargetImageHandler) GetTargetImageById(w http.ResponseWriter, r *http.Request) error {
	tih.l.Printf("Getting TargetImage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	targetImage, err := tih.repo.Get(id)
	if err != nil {
		return err
	}
	tih.l.Printf("Found TargetImage: %v", targetImage)
	response.WriteResponse(w, http.StatusOK, targetImage)
	return nil
}

func (tih *TargetImageHandler) CreateTargetImage(w http.ResponseWriter, r *http.Request) error {
	tih.l.Printf("Creating targetimage")
	var req CreateTargetImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	targetImage, err := tih.repo.Create(sqlc.CreateTargetImageParams{
		ID:          req.TargetImageId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return err
	}
	tih.l.Printf("Created target image with id: %s", targetImage.ID)
	response.WriteResponse(w, http.StatusCreated, targetImage)
	return nil
}
