package handler

import (
	"context"

	v1 "github.com/go-list-templ/proto/gen/api/sso/v1"
	pbgrpc "google.golang.org/grpc"

	"github.com/go-list-templ/sso-service/internal/port"
	"go.uber.org/zap"
)

type Auth struct {
	v1.AuthServiceServer

	authService port.AuthService
	logger      *zap.Logger
}

func RegisterAuth(s *pbgrpc.Server, a port.AuthService, l *zap.Logger) {
	service := &Auth{authService: a, logger: l}
	{
		v1.RegisterAuthServiceServer(s, service)
	}
}

func (a *Auth) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	return nil, nil
}

func (a *Auth) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	return nil, nil
}
