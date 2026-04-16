package client

import "errors"

var (
	ErrUserExists          = errors.New("exists")
	ErrUserInvalidArgument = errors.New("invalid argument")
)

type UserExists struct {
	Err error
}

func NewUserExists(err error) *UserExists {
	return &UserExists{err}
}

func (u *UserExists) Error() string {
	return u.Err.Error()
}

func (u *UserExists) Unwrap() error {
	return ErrUserExists
}

type UserInvalidArgument struct {
	Err error
}

func NewUserInvalidArgument(err error) *UserInvalidArgument {
	return &UserInvalidArgument{err}
}

func (u *UserInvalidArgument) Error() string {
	return u.Err.Error()
}

func (u *UserInvalidArgument) Unwrap() error {
	return ErrUserInvalidArgument
}
