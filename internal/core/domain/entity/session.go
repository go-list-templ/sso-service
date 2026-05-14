package entity

import (
	"time"

	"github.com/go-list-templ/sso-service/internal/core/domain/vo"
)

type Session struct {
	UserID       vo.ID
	RefreshToken vo.RefreshToken
	ExpiresAt    vo.ExpiresAt
	CreatedAt    time.Time
}

func NewSession(userID string, refreshToken string) (Session, error) {
	validRefreshToken, err := vo.NewRefreshToken(refreshToken)
	if err != nil {
		return Session{}, err
	}

	validUserId, err := vo.NewID(userID)
	if err != nil {
		return Session{}, err
	}

	now := time.Now()

	return Session{
		UserID:       validUserId,
		RefreshToken: validRefreshToken,
		ExpiresAt:    vo.NewExpiresAt(now),
		CreatedAt:    now,
	}, nil
}

func (s *Session) Expired() bool {
	return s.ExpiresAt.Expired()
}
