package handler

import (
	"context"

	v1 "github.com/go-list-templ/proto/gen/api/sso/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/internal/port"
	"go.uber.org/zap"
)

type Auth struct {
	v1.AuthServiceServer

	service port.AuthService
	logger  *zap.Logger
}

func RegisterAuth(s *pbgrpc.Server, a port.AuthService, l *zap.Logger) {
	service := &Auth{service: a, logger: l}
	{
		v1.RegisterAuthServiceServer(s, service)
	}
}

func (a *Auth) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	input := dto.AuthInput{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	output, err := a.service.Register(ctx, input)
	if err != nil {
		a.logger.Warn("register", zap.Any("context", ctx), zap.Error(err))

		return nil, err
	}

	return &v1.RegisterResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (a *Auth) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	input := dto.AuthInput{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	output, err := a.service.Register(ctx, input)
	if err != nil {
		a.logger.Warn("login", zap.Any("context", ctx), zap.Error(err))

		return nil, err
	}

	return &v1.LoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}
