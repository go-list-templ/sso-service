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

func NewSession(userID string) (Session, error) {
	id := vo.NewID()

	validUserId, err := vo.FromStr(userID)
	if err != nil {
		return Session{}, entityerr.NewSessionError("userId", err)
	}

	return Session{
		ID:           id,
		UserID:       validUserId,
		RefreshToken: "",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now(),
	}, nil
}
