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

func (ish *ImageSetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ish.l.Printf("Handling request: %s %s", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		ish.get(w, r)
	case http.MethodPost:
		ish.post(w, r)
	case http.MethodPut:
		ish.put(w, r)
	case http.MethodDelete:
		ish.delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ish *ImageSetHandler) get(w http.ResponseWriter, r *http.Request) {
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

func (ish *ImageSetHandler) post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post ImageSet\n"))
}

func (ish *ImageSetHandler) put(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Put ImageSet\n"))
}

func (ish *ImageSetHandler) delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete ImageSet\n"))
}
