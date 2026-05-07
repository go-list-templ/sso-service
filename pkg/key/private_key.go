package key

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
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

type Keys struct {
	CurrentVersion int        `json:"current_version"`
	Versions       []*Version `json:"versions"`
}

type Version struct {
	Version    int    `json:"version"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

func loadKey(cfg *config.PrivateKey) (*rsa.PrivateKey, error) {
	data, _ := os.ReadFile(cfg.File)

	var keys Keys
	err := json.Unmarshal(data, &keys)
	if err != nil {
		return nil, err
	}

	for _, version := range keys.Versions {
		if nil == version {
			continue
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(version.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("parse private key version: %v  %w", version.Version, err)
		}

		if keys.CurrentVersion == version.Version {
			return privateKey, nil
		}
	}

	return nil, fmt.Errorf("private key not found")
}
