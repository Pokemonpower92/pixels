package handler

import (
	"log/slog"
	"net/http"

	"github.com/pokemonpower92/pixels/internal/response"
)

func HealthCheck(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	response.WriteSuccessResponse(w, http.StatusOK, "ok")
}
