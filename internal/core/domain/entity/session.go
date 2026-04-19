package entity

import (
	"time"

	"github.com/go-list-templ/sso-service/internal/core/domain/entityerr"
	"github.com/go-list-templ/sso-service/internal/core/domain/vo"
	"github.com/go-list-templ/sso-service/pkg/jwt"
)

type Session struct {
	ID           vo.ID
	UserID       vo.ID
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func NewSession(userID string) (Session, error) {
	id := vo.NewID()

	validUserId, err := vo.FromStr(userID)
	if err != nil {
		return Session{}, entityerr.NewSessionError("userId", err)
	}

	refreshToken, err := jwt.CreateRefreshToken()
	if err != nil {
		return Session{}, err
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
