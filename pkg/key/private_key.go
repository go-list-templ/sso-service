package key

import (
	"crypto/rsa"
	"os"

	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type Private struct {
	*rsa.PrivateKey
}

func NewPrivate(cfg *config.PrivateKey) (*Private, error) {
	pk, err := loadKey(cfg)
	if err != nil {
		return nil, err
	}

	return &Private{pk}, nil
}

func loadKey(cfg *config.PrivateKey) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(cfg.File)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
