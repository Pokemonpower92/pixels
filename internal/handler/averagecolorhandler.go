package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/client"
	"github.com/pokemonpower92/collagegenerator/internal/imageprocessing"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/pokemonpower92/collagegenerator/internal/store"
)

type CreateAverageColorRequest struct {
	ImagesetID     uuid.UUID `json:"imageset_id"`
	AverageColorID uuid.UUID `json:"averagecolor_id"`
}

type AverageColorHandler struct {
	repo repository.ACRepo
}

func NewAverageColorHandler(repo repository.ACRepo) *AverageColorHandler {
	return &AverageColorHandler{repo: repo}
}

func (ach *AverageColorHandler) GetAverageColors(w http.ResponseWriter, _ *http.Request, l *slog.Logger) error {
	l.Info("Getting AverageColors")
	averageColors, err := ach.repo.GetAll()
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found %d AverageColors", len(averageColors)))
	response.WriteResponse(w, http.StatusOK, averageColors)
	return nil
}

func (ach *AverageColorHandler) GetAverageColorById(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Getting AverageColor by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	averageColor, err := ach.repo.Get(id)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found AverageColor: %v", averageColor))
	response.WriteResponse(w, http.StatusOK, averageColor)
	return nil
}

func (ach *AverageColorHandler) GetByImageSetId(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Getting AverageColor by ImageSet ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	averageColors, err := ach.repo.GetByResourceId(id)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found %d AverageColors", len(averageColors)))
	response.WriteResponse(w, http.StatusOK, averageColors)
	return nil
}

func (ach *AverageColorHandler) CreateAverageColor(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Creating AverageColor")
	var req CreateAverageColorRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	fileReader := client.NewFileReader("http://filestore:8081/files/", l)
	fileResponse, err := fileReader.GetFile(req.AverageColorID)
	if err != nil {
		l.Error("Error getting file", "err", err)
		return nil
	}
	image, err := store.GetRGBA(fileResponse)
	if err != nil {
		l.Error("Error converting file", "err", err)
		return err
	}
	average := imageprocessing.CalculateAverageColor(image)
	averageColor, err := ach.repo.Create(sqlc.CreateAverageColorParams{
		ID:         req.AverageColorID,
		ImagesetID: req.ImagesetID,
		FileName:   req.AverageColorID.String(),
		R:          int32(average.R),
		G:          int32(average.G),
		B:          int32(average.B),
		A:          int32(average.A),
	})
	if err != nil {
		response.WriteResponse(w, http.StatusConflict, map[string]string{
			"error": "Average color already exists for this image",
		})
		return err
	}
	l.Info(fmt.Sprintf("Created AverageColor with id: %s", averageColor.ID))
	response.WriteResponse(w, http.StatusCreated, averageColor)
	return nil
}
