package client

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/go-list-templ/sso-service/pkg/vault"
	"go.uber.org/zap"
)

const (
	SignAlgorithm = "pkcs1v15"
	Prehashed     = true

	SignPath = "transit/sign"
	KeysPath = "transit/keys"
)

var (
	ErrInvalidSignature = errors.New("invalid signature format")
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

func (v *Vault) SignJWT(ctx context.Context, unsignedToken, version string) (string, error) {
	path := filepath.Join(SignPath, v.cfg.TransitName)

	hasher := sha256.New()
	hasher.Write([]byte(unsignedToken))
	hashed := hasher.Sum(nil)

	data := map[string]any{
		"input":               base64.StdEncoding.EncodeToString(hashed),
		"key_version":         version,
		"prehashed":           Prehashed,
		"signature_algorithm": SignAlgorithm,
	}

	secret, err := v.vault.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		v.logger.Error("sign token", zap.Any("context", ctx), zap.Error(err))

		return "", err
	}

	rawSignature := fmt.Sprintf("%v", secret.Data["signature"])

	parts := strings.Split(rawSignature, ":")
	if len(parts) < 3 {
		return "", ErrInvalidSignature
	}

	decodedSig, err := base64.StdEncoding.DecodeString(parts[len(parts)-1])
	if err != nil {
		v.logger.Error("decode signature", zap.Any("context", ctx), zap.Error(err))

		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(decodedSig)

	return signature, nil
}

func (v *Vault) GetPublicKeys(ctx context.Context) (dto.PublicKeys, error) {
	path := filepath.Join(KeysPath, v.cfg.TransitName)

	secret, err := v.vault.Logical().ReadWithContext(ctx, path)
	if err != nil {
		v.logger.Error("failed to read transit key from vault", zap.Error(err))
		return dto.PublicKeys{}, err
	}

	if secret == nil || secret.Data == nil {
		v.logger.Warn("transit key data is empty")
		return dto.PublicKeys{}, nil
	}

	rawKeys, exists := secret.Data["keys"]
	if !exists {
		v.logger.Warn("no keys found in vault response")
		return dto.PublicKeys{}, nil
	}

	keysMap, ok := rawKeys.(map[string]any)
	if !ok {
		v.logger.Error("failed to cast keys to map[string]any")
		return dto.PublicKeys{}, err
	}

	publicKeys := dto.PublicKeys{
		Keys: make([]dto.PublicKey, 0),
	}

	for version, keyData := range keysMap {
		details, ok := keyData.(map[string]any)
		if !ok {
			continue
		}

		pubKeyRaw, hasPubKey := details["public_key"]
		if !hasPubKey {
			continue
		}

		publicKey := dto.PublicKey{
			Version: version,
			Key:     pubKeyRaw.(string),
		}

		publicKeys.Keys = append(publicKeys.Keys, publicKey)
	}

	return publicKeys, nil
}
