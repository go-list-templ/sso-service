package dao

import (
	"time"

	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
	"github.com/go-list-templ/sso-service/internal/core/domain/vo"
	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func FromEntity(session entity.Session) Session {
	return Session{
		ID:           session.ID.Value(),
		UserID:       session.UserID.Value(),
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
}

func (s *Session) ToEntity() entity.Session {
	return entity.Session{
		ID:           vo.UnsafeID(s.ID),
		UserID:       vo.UnsafeID(s.UserID),
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
	}
}
