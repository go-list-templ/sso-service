package key

import (
	"crypto/rsa"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrReadFile = errors.New("read file")
	ErrParseKey = errors.New("parse private key")
)

type Private struct {
	cfg    *config.PrivateKey
	logger *zap.Logger
}

func NewPrivate(cfg *config.PrivateKey, l *zap.Logger) *Private {
	return &Private{cfg, l}
}

func (p *Private) LoadKey() (*rsa.PrivateKey, error) {
	pk := filepath.Join(p.cfg.Path, p.cfg.Name)

	data, err := os.ReadFile(pk)
	if err != nil {
		p.logger.Error(ErrReadFile.Error(), zap.Error(err))

		return nil, ErrReadFile
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		p.logger.Error(ErrParseKey.Error(), zap.Error(err))

		return nil, ErrParseKey
	}

	return privateKey, nil
}
