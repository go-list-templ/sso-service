package entityerr

import (
	"errors"
	"fmt"
)

var (
	ErrUserInvalidData = errors.New("invalid session data")
	ErrSessionExpired  = errors.New("session expired")
)

type SessionErr struct {
	Field string
	Err   error
}

func NewSessionError(field string, err error) *SessionErr {
	return &SessionErr{Field: field, Err: err}
}

func (u *SessionErr) Error() string {
	return fmt.Sprintf("invalid session %s: %v", u.Field, u.Err)
}

func (u *SessionErr) Unwrap() error {
	return ErrUserInvalidData
}
