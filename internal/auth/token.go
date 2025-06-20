package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sqlc "github.com/pokemonpower92/pixels/internal/sqlc/generated"
)

type TokenManager interface {
	Generate(user *sqlc.User) (string, error)
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

func (jm *JwtManager) Generate(user *sqlc.User) (string, error) {
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
		return "", err
	}
	return s, nil
}

func (jm *JwtManager) Verify(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return &jm.PrivateKey.PublicKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if claims.Issuer != "pixels" {
		return nil, errors.New("invalid issuer")
	}
	for _, aud := range claims.Audience {
		includes := false
		if aud == "pixels" {
			break
		}
		if !includes {
			return nil, errors.New("invalid audience")
		}
	}
	return claims, nil
}
