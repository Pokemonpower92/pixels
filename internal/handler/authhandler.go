package handler

import (
	"encoding/json"
	"log"
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

type AuthHandler struct {
	l *log.Logger
}

func NewAuthHandler() *AuthHandler {
	l := log.New(log.Writer(), "", log.LstdFlags)
	return &AuthHandler{l}
}

func (ah *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) error {
	ah.l.Printf("Authenticating user")
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ah.l.Printf("Error parsing request: %s", err)
		return err
	}
	auth := auth.Authenticate(req.UserName)
	ah.l.Printf("Authenticated user: %s", req.UserName)
	response.WriteResponse(w, http.StatusOK, auth)
	return nil
}

func (ah *AuthHandler) Authorize(w http.ResponseWriter, r *http.Request) error {
	ah.l.Printf("Authorizing user")
	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ah.l.Printf("Error parsing request: %s", err)
		return err
	}
	auth := auth.Authorize(req.UserName)
	ah.l.Printf("Authorizing user: %s", req.UserName)
	response.WriteResponse(w, http.StatusOK, auth)
	return nil
}
