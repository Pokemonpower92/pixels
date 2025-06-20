package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createTestECPrivateKeyPEM() (string, *ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", nil, err
	}

	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", nil, err
	}

	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}

	pemData := pem.EncodeToMemory(pemBlock)

	return string(pemData), privateKey, nil
}

func TestGetPrivateKey(t *testing.T) {
	validPEM, expectedKey, err := createTestECPrivateKeyPEM()
	if err != nil {
		t.Fatalf("Failed to create test private key: %v", err)
	}

	testCases := []struct {
		name        string
		privateKey  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid EC private key",
			privateKey:  validPEM,
			expectError: false,
		},
		{
			name:        "invalid PEM format",
			privateKey:  "invalid-pem-data",
			expectError: true,
			errorMsg:    "failed to parse PEM block containing the private key",
		},
		{
			name:        "empty string",
			privateKey:  "",
			expectError: true,
			errorMsg:    "failed to parse PEM block containing the private key",
		},
		{
			name: "valid PEM but not EC private key",
			privateKey: `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJANLrmHDPO3qTMA0GCSqGSIb3DQEBBQUAMBQxEjAQBgNVBAMMCWxv
Y2FsaG9zdDAeFw0yMzEwMDEwMDAwMDBaFw0yNDEwMDEwMDAwMDBaMBQxEjAQBgNV
BAMMCWxvY2FsaG9zdDBcMA0GCSqGSIb3DQEBAQUAA0sAMEgCQQDTgvwjlRHZ2T1n
-----END CERTIFICATE-----`,
			expectError: true,
		},
		{
			name: "valid PEM block but invalid EC key data",
			privateKey: `-----BEGIN EC PRIVATE KEY-----
aW52YWxpZC1kYXRh
-----END EC PRIVATE KEY-----`,
			expectError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result, err := GetPrivateKey(test.privateKey)

			if test.expectError {
				if err == nil {
					t.Errorf("Test %s failed: expected error but got none", test.name)
					return
				}
				if test.errorMsg != "" && err.Error() != test.errorMsg {
					t.Errorf(
						"Test %s failed: expected error message %q, got %q",
						test.name,
						test.errorMsg,
						err.Error(),
					)
				}
				if result != nil {
					t.Errorf("Test %s failed: expected nil result on error, got %v", test.name, result)
				}
			} else {
				if err != nil {
					t.Errorf("Test %s failed: unexpected error %v", test.name, err)
					return
				}
				if result == nil {
					t.Errorf("Test %s failed: expected private key but got nil", test.name)
					return
				}
				if !result.PublicKey.Equal(&expectedKey.PublicKey) {
					t.Errorf("Test %s failed: returned key does not match expected key", test.name)
				}
			}
		})
	}
}

func TestJwtManager_Verify(t *testing.T) {
	_, privateKey, err := createTestECPrivateKeyPEM()
	if err != nil {
		t.Fatalf("Failed to create test private key: %v", err)
	}

	jwtManager := &JwtManager{
		PrivateKey: privateKey,
	}

	validClaims := jwt.RegisteredClaims{
		Issuer:    "pixels",
		Audience:  []string{"pixels"},
		Subject:   "test-user-123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	validToken, _ := jwt.NewWithClaims(jwt.SigningMethodES256, validClaims).SignedString(privateKey)

	expiredClaims := validClaims
	expiredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))
	expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodES256, expiredClaims).SignedString(privateKey)

	wrongIssuerClaims := validClaims
	wrongIssuerClaims.Issuer = "wrong-issuer"
	wrongIssuerToken, _ := jwt.NewWithClaims(jwt.SigningMethodES256, wrongIssuerClaims).SignedString(privateKey)

	wrongAudienceClaims := validClaims
	wrongAudienceClaims.Audience = []string{"wrong-audience"}
	wrongAudienceToken, _ := jwt.NewWithClaims(jwt.SigningMethodES256, wrongAudienceClaims).SignedString(privateKey)

	multiAudienceClaims := validClaims
	multiAudienceClaims.Audience = []string{"pixels", "other-service"}
	multiAudienceToken, _ := jwt.NewWithClaims(jwt.SigningMethodES256, multiAudienceClaims).SignedString(privateKey)

	wrongAlgToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims).SignedString([]byte("secret"))

	otherKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	wrongKeyToken, _ := jwt.NewWithClaims(jwt.SigningMethodES256, validClaims).SignedString(otherKey)

	testCases := []struct {
		name        string
		token       string
		expectError bool
		errorText   string
	}{
		{
			name:        "valid token",
			token:       validToken,
			expectError: false,
		},
		{
			name:        "expired token",
			token:       expiredToken,
			expectError: true,
			errorText:   "expired",
		},
		{
			name:        "wrong issuer",
			token:       wrongIssuerToken,
			expectError: true,
			errorText:   "invalid issuer",
		},
		{
			name:        "wrong audience",
			token:       wrongAudienceToken,
			expectError: true,
			errorText:   "invalid audience",
		},
		{
			name:        "multiple audiences valid",
			token:       multiAudienceToken,
			expectError: false,
		},
		{
			name:        "wrong signing algorithm",
			token:       wrongAlgToken,
			expectError: true,
			errorText:   "unexpected signing method",
		},
		{
			name:        "wrong signing key",
			token:       wrongKeyToken,
			expectError: true,
			errorText:   "verification error",
		},
		{
			name:        "malformed token",
			token:       "invalid.token.format",
			expectError: true,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			claims, err := jwtManager.Verify(test.token)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if claims != nil {
					t.Errorf("expected nil claims but got %v", claims)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if claims == nil {
					t.Errorf("expected claims but got nil")
				}
			}
		})
	}
}
