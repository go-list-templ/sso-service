package client

import (
	"context"
	"encoding/base64"
	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/go-list-templ/sso-service/pkg/vault"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
)

const (
	Algorithm = "jws"
	SignPath  = "transit/sign"
)

type Vault struct {
	cfg    *config.Vault
	vault  *vault.Vault
	logger *zap.Logger
}

func RegisterVault(c *config.Vault, v *vault.Vault, l *zap.Logger) *Vault {
	return &Vault{
		cfg:    c,
		vault:  v,
		logger: l,
	}
}

func (v *Vault) SignJWT(ctx context.Context, unsignedToken string) (string, error) {
	path := filepath.Join(SignPath, v.cfg.TransitName)

	data := map[string]any{
		"input":                base64.StdEncoding.EncodeToString([]byte(unsignedToken)),
		"marshaling_algorithm": Algorithm,
	}

	secret, err := v.vault.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		v.logger.Error("sign token", zap.Any("context", ctx), zap.Error(err))

		return "", err
	}

	rawSignature := secret.Data["signature"].(string)

	parts := strings.Split(rawSignature, ":")
	finalSignature := parts[len(parts)-1]

	return finalSignature, nil
}
