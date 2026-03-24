package jwt

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const PrivateKey = "/run/secrets/jwt_private_key"

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(PrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
