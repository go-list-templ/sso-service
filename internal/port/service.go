package port

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/core/dto"
)

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock

type (
	AuthService interface {
		Register(context.Context, dto.AuthInput) (dto.AuthOutput, error)
		Login(context.Context, dto.AuthInput) (dto.AuthOutput, error)
		Refresh(context.Context, dto.RefreshInput) (dto.AuthOutput, error)
	}

	JWKSService interface {
		Get(context.Context) (dto.PublicKeys, error)
	}
)
