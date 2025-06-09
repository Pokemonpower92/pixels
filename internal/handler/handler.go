package handler

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/imageprocessing"
	"github.com/pokemonpower92/pixels/internal/repository"
	"github.com/pokemonpower92/pixels/internal/response"
)

type ImageMetadata struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ImageHandler struct {
	repo       repository.Repository
	resolution config.ResolutionConfig
	logger     *slog.Logger
}

func NewImageHandler(repo repository.Repository, resolution config.ResolutionConfig, logger *slog.Logger) *ImageHandler {
	return &ImageHandler{
		repo:       repo,
		resolution: resolution,
		logger:     logger,
	}
}

// GetImages returns a JSON list of image metadata (without the actual image data)
func (h *ImageHandler) GetImages(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting all images")

	images, err := h.repo.GetAll()
	if err != nil {
		l.Error("Error getting images", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}

	// Convert to metadata only (no image data)
	metadata := make([]ImageMetadata, len(images))
	for i, img := range images {
		metadata[i] = ImageMetadata{
			ID:        img.ID,
			CreatedAt: img.CreatedAt.Time,
			UpdatedAt: img.UpdatedAt.Time,
		}
	}

	l.Info("Retrieved images", "count", len(metadata))
	response.WriteSuccessResponse(w, http.StatusOK, metadata)
}

// GetImage serves the raw PNG image data
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting image by id")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing UUID", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}

	image, err := h.repo.Get(id)
	if err != nil {
		l.Error("Error getting image", "error", err, "id", id)
		response.WriteErrorResponse(w, 404, err)
		return
	}

	// Serve raw PNG data
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(image.ImageData)))
	w.WriteHeader(http.StatusOK)
	w.Write(image.ImageData)

	l.Info("Served image", "id", id, "size", len(image.ImageData))
}

// CreateImage accepts raw PNG binary data, processes it, and stores the result
func (h *ImageHandler) CreateImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating image")

	if r.Header.Get("Content-Type") != "image/png" {
		l.Error("Invalid content type", "content_type", r.Header.Get("Content-Type"))
		response.WriteErrorResponse(w, 400, fmt.Errorf("only PNG images are allowed"))
		return
	}

	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("Error reading request body", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	defer r.Body.Close()

	if len(imageData) < 8 || string(imageData[:8]) != "\x89PNG\r\n\x1a\n" {
		l.Error("Invalid PNG data - missing PNG signature")
		response.WriteErrorResponse(w, 400, fmt.Errorf("invalid PNG data"))
		return
	}

	l.Info("Processing image", "original_size", len(imageData))
	imageReader := bytes.NewReader(imageData)
	sectionMap, err := imageprocessing.GetSectionColors(imageReader, l, h.resolution)
	if err != nil {
		l.Error("Error processing image sections", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}

	processedImage := imageprocessing.CreateImage(sectionMap, h.resolution)
	var buf bytes.Buffer
	if err := png.Encode(&buf, processedImage); err != nil {
		l.Error("Error encoding processed image", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}

	processedImageData := buf.Bytes()
	l.Info("Image processed", "original_size", len(imageData), "processed_size", len(processedImageData))
	image, err := h.repo.Create(processedImageData)
	if err != nil {
		l.Error("Error creating image", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}

	metadata := ImageMetadata{
		ID:        image.ID,
		CreatedAt: image.CreatedAt.Time,
		UpdatedAt: image.UpdatedAt.Time,
	}

	l.Info("Created image", "id", image.ID, "size", len(processedImageData))
	response.WriteSuccessResponse(w, http.StatusCreated, metadata)
}
