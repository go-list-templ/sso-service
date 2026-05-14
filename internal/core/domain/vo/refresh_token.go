package vo

import "errors"

var ErrEmptyToken = errors.New("empty refresh token")

type RefreshToken struct {
	value string
}

func NewRefreshToken(token string) (RefreshToken, error) {
	if token == "" {
		return RefreshToken{}, ErrEmptyToken
	}

	return RefreshToken{
		value: token,
	}, nil
}

func UnsafeRefreshToken(token string) RefreshToken {
	return RefreshToken{value: token}
}

func (r *RefreshToken) Value() string {
	return r.value
}
