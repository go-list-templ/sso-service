package client

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"strings"

	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/go-list-templ/sso-service/pkg/vault"
	"go.uber.org/zap"
)

const (
	Algorithm = "jws"

	SignPath = "transit/sign"
	KeysPath = "transit/keys"
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

func (v *Vault) SignJWT(ctx context.Context, unsignedToken string) (dto.SignJWT, error) {
	path := filepath.Join(SignPath, v.cfg.TransitName)

	data := map[string]any{
		"input":                base64.StdEncoding.EncodeToString([]byte(unsignedToken)),
		"marshaling_algorithm": Algorithm,
	}

	secret, err := v.vault.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		v.logger.Error("sign token", zap.Any("context", ctx), zap.Error(err))

		return dto.SignJWT{}, err
	}

	keyVersion := secret.Data["key_version"].(string)
	rawSignature := secret.Data["signature"].(string)

	parts := strings.Split(rawSignature, ":")
	signature := parts[len(parts)-1]

	return dto.SignJWT{
		Version:   keyVersion,
		Signature: signature,
	}, nil
}

func (v *Vault) GetPublicKeys(ctx context.Context) ([]dto.VaultPublicKey, error) {
	path := filepath.Join(KeysPath, v.cfg.TransitName)

	secret, err := v.vault.Logical().ReadWithContext(ctx, path)
	if err != nil {
		v.logger.Error("failed to read transit key from vault", zap.Error(err))
		return nil, err
	}

	if secret == nil || secret.Data == nil {
		v.logger.Warn("transit key data is empty")
		return nil, nil
	}

	rawKeys, exists := secret.Data["keys"]
	if !exists {
		v.logger.Warn("no keys found in vault response")
		return nil, nil
	}

	keysMap, ok := rawKeys.(map[string]any)
	if !ok {
		v.logger.Error("failed to cast keys to map[string]any")
		return nil, err
	}

	publicKeys := make([]dto.VaultPublicKey, 0)

	for version, keyData := range keysMap {
		details, ok := keyData.(map[string]any)
		if !ok {
			continue
		}

		pubKeyRaw, hasPubKey := details["public_key"]
		if !hasPubKey {
			continue
		}

		publicKeys = append(publicKeys, dto.VaultPublicKey{
			Version: version,
			Key:     pubKeyRaw.(string),
		})
	}

	return publicKeys, nil
}
