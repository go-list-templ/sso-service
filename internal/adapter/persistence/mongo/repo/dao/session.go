package dao

import (
	"github.com/go-list-templ/sso-service/internal/core/domain/vo"
	"time"

	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
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

// todo add logic with unsafe vo

func (s *Session) ToEntity() entity.Session {
	return entity.Session{
		ID: vo.UnsafeID(),
	}
}
