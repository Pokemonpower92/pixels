package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/auth"
	"github.com/pokemonpower92/collagegenerator/internal/response"
)

type AuthRequest struct {
	UserName string `json:"user_name"`
}

type AuthenticationResponse struct {
	Ok      bool   `json:"ok"`
	IdToken string `json:"id_token"`
}

type AuthorizationResponse struct {
	Ok bool `json:"ok"`
}

func Authenticate(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Authenticating user")
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.Error(fmt.Sprintf("Error parsing request: %s", err))
		return err
	}
	auth := auth.Authenticate(req.UserName)
	l.Info(fmt.Sprintf("Authenticated user: %s", req.UserName))
	response.WriteResponse(w, http.StatusOK, auth)
	return nil
}

func Authorize(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Authorizing user")
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.Error(fmt.Sprintf("Error parsing request: %s", err))
		return err
	}
	auth := auth.Authorize(req.UserName)
	l.Info(fmt.Sprintf("Authorized user: %s", req.UserName))
	response.WriteResponse(w, http.StatusOK, auth)
	return nil
}
