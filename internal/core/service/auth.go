package service

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/internal/port"
)

type Auth struct {
	repo       port.AuthRepo
	userClient port.UserClient
}

func NewAuth(a port.AuthRepo, u port.UserClient) *Auth {
	return &Auth{a, u}
}

func (a *Auth) Register(ctx context.Context, input dto.AuthInput) (dto.AuthOutput, error) {
	err := input.Validate()
	if err != nil {
		return dto.AuthOutput{}, err
	}

	createInput := dto.UserCreateInput{
		Email:    input.Email,
		Password: input.Password,
	}

	createOutput, err := a.userClient.Create(ctx, createInput)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	session, err := entity.NewSession(createOutput.UserId)
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
