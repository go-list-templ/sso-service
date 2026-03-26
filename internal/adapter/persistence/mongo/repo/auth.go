package repo

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/adapter/persistence/mongo"
	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
	"go.uber.org/zap"
)

type Auth struct {
	*mongo.Mongo

	logger *zap.Logger
}

func NewAuth(m *mongo.Mongo, l *zap.Logger) *Auth {
	return &Auth{m, l}
}

func (a Auth) Register(context.Context, entity.Session) error {
	return nil
}

func (a Auth) Login(context.Context, entity.Session) error {
	return nil
}
