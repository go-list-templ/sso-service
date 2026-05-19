package interceptor

import (
	"context"
	"errors"

	"github.com/go-list-templ/sso-service/internal/adapter/grpc/client"
	"github.com/go-list-templ/sso-service/internal/core/domain/entityerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ErrInternalServer = "internal server"

var allErr = map[error]codes.Code{
	client.ErrUserExists:          codes.AlreadyExists,
	client.ErrUserInvalidArgument: codes.InvalidArgument,
	entityerr.ErrSessionExpired:   codes.Unauthenticated,
	entityerr.ErrSessionNotFound:  codes.NotFound,
}

func ErrorHandling() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, toGrpcError(err)
		}

		return resp, nil
	}
}

func toGrpcError(err error) error {
	for currentErr, resCode := range allErr {
		if errors.Is(err, currentErr) {
			return status.Error(resCode, err.Error())
		}
	}

	return status.Error(codes.Internal, ErrInternalServer)
}
