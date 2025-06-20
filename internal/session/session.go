package session

import (
	"github.com/google/uuid"
	"github.com/pokemonpower92/pixels/internal/auth"
	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"
)

type Sessionizer interface {
	CreateSession(uuid.UUID) (string, error)
	FindSession(string) (uuid.UUID, bool)
}

type JWTSessionizer struct {
	jwtManager *auth.JwtManager
}

func NewJWTSessionizer(jwtManager *auth.JwtManager) *JWTSessionizer {
	return &JWTSessionizer{jwtManager: jwtManager}
}

func (js *JWTSessionizer) CreateSession(userID uuid.UUID) (string, error) {
	user := &sqlc.User{ID: userID}
	return js.jwtManager.Generate(user)
}

func (js *JWTSessionizer) FindSession(token string) (uuid.UUID, bool) {
	claims, err := js.jwtManager.Verify(token)
	if err != nil {
		return uuid.UUID{}, false
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, false
	}
	return userID, true
}
