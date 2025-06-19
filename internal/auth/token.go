package jwt

import (
	"crypto/ecdsa"

	"github.com/golang-jwt/jwt/v5"
	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"
)

type TokenManager interface {
	Generate(user *sqlc.User) (error, string)
}

type JwtManager struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (jm *JwtManager) Generate(user *sqlc.User) (error, string) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss": "pixels",
			"sub": user.ID.String(),
			"aud": "pixels",
		},
	)
	s, err := t.SignedString(jm.privateKey)
	if err != nil {
		return err, ""
	}
	return nil, s
}

func (jm *JwtManager) Verify(token string) error {
	return nil
}
