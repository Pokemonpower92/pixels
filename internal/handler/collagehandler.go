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
	l := log.New(log.Writer(), "", log.LstdFlags)
	return &CollageHandler{l: l, repo: repo}
}

func (ch *CollageHandler) GetCollages(w http.ResponseWriter, _ *http.Request) error {
	ch.l.Printf("Getting Collages")
	imageSets, err := ch.repo.GetAll()
	if err != nil {
		return err
	}
	ch.l.Printf("Found %d Collages", len(imageSets))
	response.WriteResponse(w, http.StatusOK, imageSets)
	return nil
}

func (ch *CollageHandler) GetCollageById(w http.ResponseWriter, r *http.Request) error {
	ch.l.Printf("Getting Collage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	averageColor, err := ch.repo.Get(id)
	if err != nil {
		return err
	}
	ch.l.Printf("Found Collage: %v", averageColor)
	response.WriteResponse(w, http.StatusOK, averageColor)
	return nil
}

func (ch *CollageHandler) CreateCollage(w http.ResponseWriter, r *http.Request) error {
	ch.l.Printf("Creating Collage")
	var req sqlc.CreateCollageParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	collage, err := ch.repo.Create(req)
	if err != nil {
		return err
	}
	ch.l.Printf("Created Collage with id: %s", collage.ID)
	go service.CreateCollageMetaData(collage)
	response.WriteResponse(w, http.StatusCreated, collage)
	return nil
}
