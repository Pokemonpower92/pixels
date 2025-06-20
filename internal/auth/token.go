package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"
)

type TokenManager interface {
	Generate(user *sqlc.User) (error, string)
}

type JwtManager struct {
	PrivateKey *ecdsa.PrivateKey
}

func GetPrivateKey(privateKeyPEM string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

func (jm *JwtManager) Generate(user *sqlc.User) (error, string) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		jwt.RegisteredClaims{
			Issuer:    "pixels",
			Audience:  []string{"pixels"},
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	)
	s, err := t.SignedString(jm.PrivateKey)
	if err != nil {
		return err, ""
	}
	return nil, s
}

func (jm *JwtManager) Verify(token string) error {
	return nil
}
