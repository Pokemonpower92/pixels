package handler

import (
	"bytes"
	"errors"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pokemonpower92/pixels/config"
	"github.com/pokemonpower92/pixels/internal/imageprocessing"
	"github.com/pokemonpower92/pixels/internal/middleware"
	"github.com/pokemonpower92/pixels/internal/repository"
	"github.com/pokemonpower92/pixels/internal/response"
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
	repo       repository.ImageModeler
	resolution config.ResolutionConfig
	logger     *slog.Logger
}

func NewImageHandler(
	repo repository.ImageModeler,
	resolution config.ResolutionConfig,
	logger *slog.Logger,
) *ImageHandler {
	return &ImageHandler{
		repo:       repo,
		resolution: resolution,
		logger:     logger,
	}
}

// GetImages returns a JSON list of image metadata (without the actual image data)
func (h *ImageHandler) GetImages(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting all images")
	userId, ok := middleware.GetUserID(r)
	if !ok {
		l.Error("Error fetching session information")
		response.WriteErrorResponse(w, 401, errors.New("Could not find user session"))
	}
	images, err := h.repo.GetAll(uuid.MustParse(userId))
	if err != nil {
		l.Error("Error getting images", "error", err)
		response.WriteErrorResponse(w, 500, err)
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
	response.WriteSuccessResponse(w, http.StatusOK, metadata)
}

// GetImage serves the raw PNG image data
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting image by id")
	userId, ok := middleware.GetUserID(r)
	if !ok {
		l.Error("Error fetching session information")
		response.WriteErrorResponse(w, 401, errors.New("Could not find user session"))
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing UUID", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	image, err := h.repo.Get(sqlc.GetImageParams{
		UserID: uuid.MustParse(userId),
		ID:     id,
	})
	if err != nil {
		l.Error("Error getting image", "error", err, "id", id)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	response.WriteImageResponse(w, image)
	l.Info("Served image", "id", id, "size", len(image.ImageData))
}

// CreateImage accepts raw PNG binary data, processes it, and stores the result
func (h *ImageHandler) CreateImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating image")
	userId, ok := middleware.GetUserID(r)
	if !ok {
		l.Error("Error fetching session information")
		response.WriteErrorResponse(w, 401, errors.New("Could not find user session"))
	}
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("Error reading request body", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	defer r.Body.Close()

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
	image, err := h.repo.Create(
		sqlc.CreateImageParams{
			UserID:    uuid.MustParse(userId),
			ImageData: processedImageData,
		},
	)
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
	response.WriteSuccessResponse(w, http.StatusCreated, metadata)
	l.Info("Created image", "id", image.ID, "size", len(processedImageData))
}
