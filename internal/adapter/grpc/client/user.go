package client

import (
	"context"

	v1 "github.com/go-list-templ/proto/gen/api/user/v1"

	"github.com/go-list-templ/sso-service/internal/core/dto"
	"github.com/go-list-templ/sso-service/pkg/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type User struct {
	client v1.UserServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

func RegisterUser(cfg *config.UserClient, l *zap.Logger) (*User, error) {
	grpcConn, err := grpc.NewClient(
		cfg.Host+":"+cfg.Port,
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &User{
		client: v1.NewUserServiceClient(grpcConn),
		conn:   grpcConn,
		logger: l,
	}, nil
}

func (u *User) Create(ctx context.Context, input dto.UserCreateInput) (dto.UserCreateOutput, error) {
	request := &v1.CreateRequest{
		Email:    input.Email,
		Password: input.Password,
	}

	response, err := u.client.Create(ctx, request)
	if err != nil {
		u.logger.Error("create user", zap.Any("context", ctx), zap.Error(err))

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return dto.UserCreateOutput{}, NewUserInvalidArgument(st.Message(), err)
			case codes.AlreadyExists:
				return dto.UserCreateOutput{}, NewUserExists(st.Message(), err)
			default:
				return dto.UserCreateOutput{}, err
			}
		}
	}

	return dto.UserCreateOutput{
		UserId: response.User.Id,
	}, err
}

func (u *User) VerifyCred(ctx context.Context, input dto.UserVerifyCredInput) (dto.UserVerifyCredOutput, error) {
	request := &v1.VerifyCredRequest{
		Email:    input.Email,
		Password: input.Password,
	}

	response, err := u.client.VerifyCred(ctx, request)
	if err != nil {
		u.logger.Error("verify cred user", zap.Any("context", ctx), zap.Error(err))

		return dto.UserVerifyCredOutput{}, err
	}

	return dto.UserVerifyCredOutput{
		UserId: response.UserId,
	}, err
}
