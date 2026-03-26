package service

import (
	"context"
	"github.com/go-list-templ/sso-service/internal/core/domain/entity"

	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/internal/port"
)

type Auth struct {
	repo port.AuthRepo
}

func NewAuth(a port.AuthRepo) *Auth {
	return &Auth{a}
}

func (a *Auth) Register(ctx context.Context, input dto.AuthInput) (dto.AuthOutput, error) {
	err := input.Validate()
	if err != nil {
		return dto.AuthOutput{}, err
	}

	//todo create user from users-service

	session, err := entity.NewSession("")
	if err != nil {
		return dto.AuthOutput{}, err
	}

	if err = a.repo.Register(ctx, session); err != nil {
		return dto.AuthOutput{}, err
	}

	return dto.AuthOutput{
		AccessToken:  "",
		RefreshToken: session.RefreshToken,
	}, nil
}

func (a *Auth) Login(ctx context.Context, input dto.AuthInput) (dto.AuthOutput, error) {
	//todo add verify email from users-service

	return dto.AuthOutput{}, nil
}
