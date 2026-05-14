package dao

import (
	"time"

	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
	"github.com/go-list-templ/sso-service/internal/core/domain/vo"
	"github.com/google/uuid"
)

type Session struct {
	UserID       uuid.UUID `bson:"user_id"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiresAt    time.Time `bson:"expires_at"`
	CreatedAt    time.Time `bson:"created_at"`
}

func FromEntity(session entity.Session) Session {
	return Session{
		UserID:       session.UserID.Value(),
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}

func (s *Session) ToEntity() entity.Session {
	return entity.Session{
		UserID:       vo.UnsafeID(s.UserID),
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
	}
}
