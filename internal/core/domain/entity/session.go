package entity

import (
	"time"

	"github.com/go-list-templ/sso-service/internal/core/domain/entityerr"
	"github.com/go-list-templ/sso-service/internal/core/domain/vo"
)

type Session struct {
	ID           vo.ID
	UserID       vo.ID
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func NewSession(userID string, refreshToken string) (Session, error) {
	id := vo.NewID()

	validUserId, err := vo.FromStr(userID)
	if err != nil {
		return Session{}, entityerr.NewSessionError("userId", err)
	}

	now := time.Now()

	return Session{
		ID:           id,
		UserID:       validUserId,
		RefreshToken: refreshToken,
		CreatedAt:    now,
		ExpiresAt:    now,
	}, nil
}
