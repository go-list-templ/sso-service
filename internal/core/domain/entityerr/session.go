package entityerr

import (
	"errors"
	"fmt"
)

var (
	ErrUserInvalidData = errors.New("invalid session data")
)

type UserError struct {
	Field string
	Err   error
}

func NewSessionError(field string, err error) *UserError {
	return &UserError{Field: field, Err: err}
}

func (u *UserError) Error() string {
	return fmt.Sprintf("invalid session %s: %v", u.Field, u.Err)
}

func (u *UserError) Unwrap() error {
	return ErrUserInvalidData
}
