package service

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/internal/port"
)

type Auth struct {
	authRepo port.AuthRepo
}

func NewAuth() {}

func (a *Auth) Register(context.Context, dto.AuthInput) (dto.AuthOutput, error) {
	return dto.AuthOutput{}, nil
}

func (a *Auth) Login(context.Context, dto.AuthInput) (dto.AuthOutput, error) {
	return dto.AuthOutput{}, nil
}
