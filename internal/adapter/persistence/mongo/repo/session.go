package repo

import (
	"context"
	"github.com/go-list-templ/sso-service/internal/adapter/persistence/mongo"
	"github.com/go-list-templ/sso-service/internal/adapter/persistence/mongo/repo/dao"
	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

const Collection = "sessions"

type Session struct {
	*mongo.Mongo

	logger *zap.Logger
}

func NewSession(m *mongo.Mongo, l *zap.Logger) *Session {
	return &Session{m, l}
}

//TODO auto delete in mongo after refresh token expired

func (s Session) Store(ctx context.Context, session entity.Session) error {
	collection := s.Database.Collection(Collection)

	sessionDAO := dao.FromEntity(session)

	_, err := collection.InsertOne(ctx, sessionDAO)
	if err != nil {
		s.logger.Error("insert session", zap.Any("ctx", ctx), zap.Error(err))
	}

	return err
}

func (s Session) FindAndDelete(ctx context.Context, refreshToken string) (entity.Session, error) {
	collection := s.Database.Collection(Collection)

	filter := bson.M{"refresh_token": refreshToken}

	var sessionDAO dao.Session

	err := collection.FindOneAndDelete(ctx, filter).Decode(&sessionDAO)
	if err != nil {
		s.logger.Error("delete session", zap.Any("ctx", ctx), zap.Error(err))

		return entity.Session{}, err
	}

	return sessionDAO.ToEntity(), nil
}
