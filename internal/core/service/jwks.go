package service

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/internal/port"
)

type JWKS struct {
	vaultClient port.VaultClient
}

func NewJWKS(v port.VaultClient) *JWKS {
	return &JWKS{v}
}

func (a *JWKS) Get(ctx context.Context) ([]dto.VaultPublicKey, error) {
	return a.vaultClient.GetPublicKeys(ctx)
}
