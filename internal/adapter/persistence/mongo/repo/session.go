package repo

import (
	"context"
	"errors"

	mongoClient "github.com/go-list-templ/sso-service/internal/adapter/persistence/mongo"

	"github.com/go-list-templ/sso-service/internal/adapter/persistence/mongo/repo/dao"
	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
	"github.com/go-list-templ/sso-service/internal/core/domain/entityerr"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

const Collection = "sessions"

type Session struct {
	*mongoClient.Mongo

	logger *zap.Logger
}

func NewSession(m *mongoClient.Mongo, l *zap.Logger) *Session {
	return &Session{m, l}
}

//TODO auto delete in mongo after refresh token expired

func (s *Session) Store(ctx context.Context, session entity.Session) error {
	collection := s.Database.Collection(Collection)

	sessionDAO := dao.FromEntity(session)

	_, err := collection.InsertOne(ctx, sessionDAO)
	if err != nil {
		s.logger.Error("store", zap.Any("context", ctx), zap.Error(err))
	}

	return s.toMongoError(err)
}

func (s *Session) FindAndDelete(ctx context.Context, refreshToken string) (entity.Session, error) {
	collection := s.Database.Collection(Collection)

	filter := bson.M{"refresh_token": refreshToken}

	var sessionDAO dao.Session

	err := collection.FindOneAndDelete(ctx, filter).Decode(&sessionDAO)
	if err != nil {
		s.logger.Error("find and delete", zap.Any("context", ctx), zap.Error(err))

		return entity.Session{}, s.toMongoError(err)
	}

	return sessionDAO.ToEntity(), nil
}

func (s *Session) toMongoError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return entityerr.ErrSessionNotFound
	}

	return err
}
