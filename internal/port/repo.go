package port

import (
	"context"
	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
)

//go:generate mockgen -source=repo.go -destination=mock/mock_repo.go -package=mock

type (
	SessionRepo interface {
		Store(context.Context, entity.Session) error
	}
)
