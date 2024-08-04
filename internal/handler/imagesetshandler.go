package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type ImageSetsHandler struct {
	l    *log.Logger
	repo isRepository
}

func NewImageSetsHandler(l *log.Logger, repo isRepository) *ImageSetsHandler {
	return &ImageSetsHandler{l: l, repo: repo}
}

func (ish *ImageSetsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (ish *ImageSetsHandler) get(w http.ResponseWriter, r *http.Request) {
	ish.l.Printf("Getting ImageSets")
	imageSets, ok := ish.repo.GetAll()
	if !ok {
		ish.l.Printf("Failed to get ImageSets")
		http.Error(w, "Error getting ImageSets", http.StatusInternalServerError)
		return
	}
	ish.l.Printf("Found %d ImageSets.", len(imageSets))

	encoder := json.NewEncoder(w)
	err := encoder.Encode(imageSets)
	if err != nil {
		ish.l.Printf("Failed to encode ImageSets: %s", err)
		http.Error(w, "Error encoding ImageSets", http.StatusInternalServerError)
	}
}

func (ish *ImageSetsHandler) post(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (ish *ImageSetsHandler) put(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (ish *ImageSetsHandler) delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
