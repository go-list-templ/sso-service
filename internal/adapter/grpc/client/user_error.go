package client

import "errors"

var (
	ErrUserExists          = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserInvalidArgument = errors.New("user invalid argument")
)

type UserExists struct {
	Message string
	Err     error
}

func NewUserExists(m string, err error) *UserExists {
	return &UserExists{m, err}
}

func (u *UserExists) Error() string {
	return u.Message
}

func (u *UserExists) Unwrap() error {
	return ErrUserExists
}

type UserInvalidArgument struct {
	Message string
	Err     error
}

func NewUserInvalidArgument(m string, err error) *UserInvalidArgument {
	return &UserInvalidArgument{m, err}
}

func (u *UserInvalidArgument) Error() string {
	return u.Message
}

func (u *UserInvalidArgument) Unwrap() error {
	return ErrUserInvalidArgument
}

type UserNotFound struct {
	Message string
	Err     error
}

func NewUserNotFound(m string, err error) *UserNotFound {
	return &UserNotFound{m, err}
}

func (u *UserNotFound) Error() string {
	return u.Message
}

func (u *UserNotFound) Unwrap() error {
	return ErrUserNotFound
}
