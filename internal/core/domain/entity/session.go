package entity

import (
	"github.com/go-list-templ/sso-service/internal/core/domain/entityerr"
	"time"

	"github.com/go-list-templ/sso-service/internal/core/domain/vo"
)

type Session struct {
	ID        vo.ID
	UserID    vo.ID
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewSession(userID string) (Session, error) {
	id := vo.NewID()

	validUserId, err := vo.FromStr(userID)
	if err != nil {
		return Session{}, entityerr.NewSessionError("userId", err)
	}

	return Session{
		ID:        id,
		UserID:    validUserId,
		Token:     "",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
	}, nil
}
