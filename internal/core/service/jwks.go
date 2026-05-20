package service

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/internal/port"
	"github.com/go-list-templ/sso-service/pkg/jwks"
)

type JWKS struct {
	vaultClient port.VaultClient
	jwks        *jwks.JWKS
}

func NewJWKS(v port.VaultClient, j *jwks.JWKS) *JWKS {
	return &JWKS{v, j}
}

func (a *JWKS) Get(ctx context.Context) (dto.JWKS, error) {
	pk, err := a.vaultClient.GetPublicKeys(ctx)
	if err != nil {
		return dto.JWKS{}, err
	}

	return a.jwks.FromPublicKeys(pk)
}
