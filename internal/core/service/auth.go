package service

import (
	"context"

	"github.com/go-list-templ/sso-service/internal/core/domain/entity"
	"github.com/go-list-templ/sso-service/internal/core/domain/entityerr"
	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/internal/port"
	"github.com/go-list-templ/sso-service/pkg/token"
)

type Auth struct {
	userClient  port.UserClient
	vaultClient port.VaultClient
	repo        port.SessionRepo
	token       *token.Token
}

func NewAuth(u port.UserClient, v port.VaultClient, r port.SessionRepo, t *token.Token) *Auth {
	return &Auth{u, v, r, t}
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

	return a.createSession(ctx, createOutput.UserId)
}

func (a *Auth) Login(ctx context.Context, input dto.AuthInput) (dto.AuthOutput, error) {
	verifyCredInput := dto.UserVerifyCredInput{
		Email:    input.Email,
		Password: input.Password,
	}

	verifyCredOutput, err := a.userClient.VerifyCred(ctx, verifyCredInput)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	return a.createSession(ctx, verifyCredOutput.UserId)
}

func (a *Auth) Refresh(ctx context.Context, input dto.RefreshInput) (dto.AuthOutput, error) {
	currentSession, err := a.repo.FindAndDelete(ctx, input.RefreshToken)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	if currentSession.Expired() {
		return dto.AuthOutput{}, entityerr.ErrSessionExpired
	}

	return a.createSession(ctx, currentSession.UserID.Value().String())
}

func (a *Auth) createSession(ctx context.Context, userId string) (dto.AuthOutput, error) {
	refreshToken, err := a.token.CreateRefresh()
	if err != nil {
		return dto.AuthOutput{}, err
	}

	session, err := entity.NewSession(userId, refreshToken)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	if err = a.repo.Store(ctx, session); err != nil {
		return dto.AuthOutput{}, err
	}

	unsignedToken, err := a.token.Unsigned()
	if err != nil {
		return dto.AuthOutput{}, err
	}

	signature, err := a.vaultClient.SignJWT(ctx, unsignedToken)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	accessToken, err := a.token.CreateAccess(unsignedToken, signature)
	if err != nil {
		return dto.AuthOutput{}, err
	}

	return dto.AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}
