package port

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/core/dto"
)

//go:generate mockgen -source=client.go -destination=mock/mock_client.go -package=mock

type (
	UserClient interface {
		Create(context.Context, dto.UserCreateInput) (dto.UserCreateOutput, error)
		VerifyCred(context.Context, dto.UserVerifyCredInput) (dto.UserVerifyCredOutput, error)
	}
)
