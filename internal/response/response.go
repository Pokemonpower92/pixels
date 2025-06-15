package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"
)

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
