package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/pokemonpower92/collagegenerator/internal/utils"
)

type AverageColorHandler struct {
	l     *log.Logger
	repo  repository.ACRepo
	store datastore.Store
}

func NewAverageColorHandler(repo repository.ACRepo) *AverageColorHandler {
	l := log.New(log.Writer(), "", log.LstdFlags)
	store := datastore.NewStore()
	return &AverageColorHandler{l: l, repo: repo, store: store}
}

func (ach *AverageColorHandler) GetAverageColors(w http.ResponseWriter, _ *http.Request) error {
	ach.l.Printf("Getting AverageColors")
	imageSets, err := ach.repo.GetAll()
	if err != nil {
		return err
	}
	ach.l.Printf("Found %d AverageColors", len(imageSets))
	response.WriteResponse(w, http.StatusOK, imageSets)
	return nil
}

func (ach *AverageColorHandler) GetAverageColorById(w http.ResponseWriter, r *http.Request) error {
	ach.l.Printf("Getting AverageColor by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	averageColor, err := ach.repo.Get(id)
	if err != nil {
		return err
	}
	ach.l.Printf("Found AverageColor: %v", averageColor)
	response.WriteResponse(w, http.StatusOK, averageColor)
	return nil
}

func (ach *AverageColorHandler) CreateAverageColor(w http.ResponseWriter, r *http.Request) error {
	ach.l.Printf("Creating AverageColor")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}
	file, handle, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(fileBytes)
	image, err := utils.ImageFileToRGBA(reader)
	if err != nil {
		return err
	}
	average := utils.CalculateAverageColor(image)
	imageSetId := uuid.MustParse(r.FormValue("imageSetId"))
	dbParams := sqlc.CreateAverageColorParams{
		ImagesetID: imageSetId,
		FileName:   handle.Filename,
		R:          int32(average.R),
		G:          int32(average.G),
		B:          int32(average.B),
		A:          int32(average.A),
	}
	averageColor, err := ach.repo.Create(dbParams)
	if err != nil {
		return err
	}
	if err := ach.store.PutImage(averageColor.ID, reader); err != nil {
		return err
	}
	ach.l.Printf("Created AverageColor with id: %s", averageColor.ID.String())
	response.WriteResponse(w, http.StatusCreated, averageColor)
	return nil
}
