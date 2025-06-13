package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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
	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
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

func NewImageHandler(
	repo repository.Repository,
	resolution config.ResolutionConfig,
	logger *slog.Logger,
) *ImageHandler {
	return &ImageHandler{
		repo:       repo,
		resolution: resolution,
		logger:     logger,
	}
}

type SuccessResponse struct {
	StatusCode int `json:"status_code"`
	Data       any `json:"data"`
}

type ErrorResponse struct {
	StatusCode int `json:"status_code"`
	Error      any `json:"error"`
}

// WriteSuccessResponse writes a successful response to w
func WriteSuccessResponse(w http.ResponseWriter, status int, val any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&SuccessResponse{
		StatusCode: status,
		Data:       val,
	})
}

// WriteErrorResponse writes an error to w.
func WriteErrorResponse(w http.ResponseWriter, status int, err error) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&ErrorResponse{
		StatusCode: status,
		Error:      err,
	})
}

// WriteImageResponse writes an image to w
func WriteImageResponse(w http.ResponseWriter, image *sqlc.Image) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(image.ImageData)))
	w.WriteHeader(http.StatusOK)
	w.Write(image.ImageData)
}

// GetImages returns a JSON list of image metadata (without the actual image data)
func (h *ImageHandler) GetImages(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting all images")
	images, err := h.repo.GetAll()
	if err != nil {
		l.Error("Error getting images", "error", err)
		WriteErrorResponse(w, 500, err)
		return
	}
	metadata := make([]ImageMetadata, len(images))
	for i, img := range images {
		metadata[i] = ImageMetadata{
			ID:        img.ID,
			CreatedAt: img.CreatedAt.Time,
			UpdatedAt: img.UpdatedAt.Time,
		}
	}
	l.Info("Retrieved images", "count", len(metadata))
	WriteSuccessResponse(w, http.StatusOK, metadata)
}

// GetImage serves the raw PNG image data
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting image by id")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing UUID", "error", err)
		WriteErrorResponse(w, 422, err)
		return
	}
	image, err := h.repo.Get(id)
	if err != nil {
		l.Error("Error getting image", "error", err, "id", id)
		WriteErrorResponse(w, 404, err)
		return
	}
	WriteImageResponse(w, image)
	l.Info("Served image", "id", id, "size", len(image.ImageData))
}

func validateCreate(r *http.Request) error {
	switch contentType := r.Header.Get("Content-Type"); contentType {
	case "image/png":
	case "image/jpeg":
	case "image/webp":
		return nil
	default:
		// return errors.New(fmt.Sprintf("Unsupported Content-Type: %s", contentType))
		return nil
	}
	return errors.New("Could not validate create request")
}

// CreateImage accepts raw PNG binary data, processes it, and stores the result
func (h *ImageHandler) CreateImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating image")
	if err := validateCreate(r); err != nil {
		l.Error("Invalid content type", "content_type", r.Header.Get("Content-Type"))
		WriteErrorResponse(w, 400, fmt.Errorf("only PNG images are allowed"))
		return
	}
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("Error reading request body", "error", err)
		WriteErrorResponse(w, 500, err)
		return
	}
	defer r.Body.Close()

	l.Info("Processing image", "original_size", len(imageData))
	imageReader := bytes.NewReader(imageData)
	sectionMap, err := imageprocessing.GetSectionColors(imageReader, l, h.resolution)
	if err != nil {
		l.Error("Error processing image sections", "error", err)
		WriteErrorResponse(w, 500, err)
		return
	}
	processedImage := imageprocessing.CreateImage(sectionMap, h.resolution)
	var buf bytes.Buffer
	if err := png.Encode(&buf, processedImage); err != nil {
		l.Error("Error encoding processed image", "error", err)
		WriteErrorResponse(w, 500, err)
		return
	}
	processedImageData := buf.Bytes()
	l.Info("Image processed", "original_size", len(imageData), "processed_size", len(processedImageData))
	image, err := h.repo.Create(processedImageData)
	if err != nil {
		l.Error("Error creating image", "error", err)
		WriteErrorResponse(w, 500, err)
		return
	}
	metadata := ImageMetadata{
		ID:        image.ID,
		CreatedAt: image.CreatedAt.Time,
		UpdatedAt: image.UpdatedAt.Time,
	}
	WriteSuccessResponse(w, http.StatusCreated, metadata)
	l.Info("Created image", "id", image.ID, "size", len(processedImageData))
}

func HealthCheck(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	WriteSuccessResponse(w, http.StatusOK, "ok")
}
