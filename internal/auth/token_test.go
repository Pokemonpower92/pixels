package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func createTestECPrivateKeyPEM() (string, *ecdsa.PrivateKey, error) {
	// Generate a new ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", nil, err
	}

	// Marshal the private key to DER format
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", nil, err
	}

	// Create PEM block
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}

	// Encode to PEM format
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
				// Verify the key is the same by comparing the public key
				if !result.PublicKey.Equal(&expectedKey.PublicKey) {
					t.Errorf("Test %s failed: returned key does not match expected key", test.name)
				}
			}
		})
	}
}
