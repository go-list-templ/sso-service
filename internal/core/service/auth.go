package service

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/internal/port"
	"github.com/go-list-templ/sso-service/pkg/token"
)

type Auth struct {
	repo       port.SessionRepo
	userClient port.UserClient
	token      *token.Token
}

func NewAuth(a port.SessionRepo, u port.UserClient, t *token.Token) *Auth {
	return &Auth{a, u, t}
}

func (a *Auth) Register(ctx context.Context, input dto.AuthInput) (dto.AuthOutput, error) {
	createInput := dto.UserCreateInput{
		Email:    input.Email,
		Password: input.Password,
	}

	createOutput, err := a.userClient.Create(ctx, createInput)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	refreshToken, err := a.token.CreateRefresh()
	if err != nil {
		return dto.AuthOutput{}, err
	}

	session, err := entity.NewSession(createOutput.UserId, refreshToken)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	if err = a.repo.Store(ctx, session); err != nil {
		return dto.AuthOutput{}, err
	}

	accessToken, err := a.token.CreateAccess(ctx)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	return dto.AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (a *Auth) Login(ctx context.Context, input dto.AuthInput) (dto.AuthOutput, error) {
	//todo add verify email from users-service

	return dto.AuthOutput{}, nil
}
