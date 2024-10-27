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

type AverageColorHandler struct {
	l    *log.Logger
	repo repository.ACRepo
}

func NewAverageColorHandler(repo repository.ACRepo) *AverageColorHandler {
	l := log.New(log.Writer(), "AverageColorHandler: ", log.LstdFlags)
	return &AverageColorHandler{l: l, repo: repo}
}

func (acr *AverageColorHandler) GetAverageColors(w http.ResponseWriter, _ *http.Request) error {
	acr.l.Printf("Getting AverageColors")
	imageSets, err := acr.repo.GetAll()
	if err != nil {
		return err
	}
	acr.l.Printf("Found %d AverageColors", len(imageSets))
	response.WriteResponse(w, http.StatusOK, imageSets)
	return nil
}

func (acr *AverageColorHandler) GetAverageColorById(w http.ResponseWriter, r *http.Request) error {
	acr.l.Printf("Getting AverageColor by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	averageColor, err := acr.repo.Get(id)
	if err != nil {
		return err
	}
	acr.l.Printf("Found AverageColor: %v", averageColor)
	response.WriteResponse(w, http.StatusOK, averageColor)
	return nil
}

func (acr *AverageColorHandler) CreateAverageColor(w http.ResponseWriter, r *http.Request) error {
	acr.l.Printf("Creating AverageColor.")
	var req sqlc.CreateAverageColorParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	averageColor, err := acr.repo.Create(req)
	if err != nil {
		return err
	}
	acr.l.Printf("Created AverageColor with id: %s", averageColor.ID)
	response.WriteResponse(w, http.StatusCreated, averageColor)
	return nil
}
