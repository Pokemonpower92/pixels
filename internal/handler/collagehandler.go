package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	"github.com/pokemonpower92/collagegenerator/internal/service"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageHandler struct {
	l    *log.Logger
	repo repository.CRepo
}

func NewCollageHandler(repo repository.CRepo) *CollageHandler {
	l := log.New(log.Writer(), "CollageHandler: ", log.LstdFlags)
	return &CollageHandler{l: l, repo: repo}
}

func (acr *CollageHandler) GetCollages(w http.ResponseWriter, _ *http.Request) error {
	acr.l.Printf("Getting Collages")
	imageSets, err := acr.repo.GetAll()
	if err != nil {
		return err
	}
	acr.l.Printf("Found %d Collages", len(imageSets))
	response.WriteResponse(w, http.StatusOK, imageSets)
	return nil
}

func (acr *CollageHandler) GetCollageById(w http.ResponseWriter, r *http.Request) error {
	acr.l.Printf("Getting Collage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	averageColor, err := acr.repo.Get(id)
	if err != nil {
		return err
	}
	acr.l.Printf("Found Collage: %v", averageColor)
	response.WriteResponse(w, http.StatusOK, averageColor)
	return nil
}

func (acr *CollageHandler) CreateCollage(w http.ResponseWriter, r *http.Request) error {
	acr.l.Printf("Creating Collage.")
	var req sqlc.CreateCollageParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	collage, err := acr.repo.Create(req)
	if err != nil {
		return err
	}
	acr.l.Printf("Created Collage with id: %s", collage.ID)
	go service.CreateCollage(collage)
	response.WriteResponse(w, http.StatusCreated, collage)
	return nil
}
