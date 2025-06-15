package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pokemonpower92/pixels/internal/repository"
	"github.com/pokemonpower92/pixels/internal/response"
	"github.com/pokemonpower92/pixels/internal/session"
	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"
	"golang.org/x/crypto/bcrypt"
)

type UserMetadata struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type AuthHandler struct {
	repo        repository.UserModeler
	sessionizer session.Sessionizer
	logger      *slog.Logger
}

func NewAuthHandler(
	repo repository.UserModeler,
	sessionizer session.Sessionizer,
	logger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		repo:        repo,
		sessionizer: sessionizer,
		logger:      logger,
	}
}

// CreateUser creates a user
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating user")
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.Error("Error decoding CreateUser request", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Error("Error hashing password", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	query := sqlc.CreateUserParams{
		UserName: req.UserName,
		Password: string(hashedPassword),
	}
	user, err := h.repo.Create(query)
	if err != nil {
		l.Error("Error creating user", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	response.WriteSuccessResponse(
		w,
		200,
		UserMetadata{
			ID:        user.ID,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
		},
	)
}

// Login logs a user in.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Logging in user")
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.Error("Error decoding Login request", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	user, err := h.repo.Get(req.UserName)
	if err != nil {
		l.Error("Error getting User", "user", req.UserName, "error", err)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	unAuthenticated := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if unAuthenticated != nil {
		l.Error("Failed login attemp for User", "user", user.ID, "error", err)
		response.WriteErrorResponse(w, 403, err)
		return
	}
	sessionId := h.sessionizer.CreateSession(user.ID)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
	})
	response.WriteSuccessResponse(
		w,
		200,
		map[string]string{
			"status": "ok",
		},
	)
}
