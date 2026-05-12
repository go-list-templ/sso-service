package interceptor

import (
	"context"
	"errors"

	"github.com/go-list-templ/sso-service/internal/adapter/grpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ErrInternalServer = "internal server"

var allErr = map[error]codes.Code{
	client.ErrUserExists:          codes.AlreadyExists,
	client.ErrUserInvalidArgument: codes.InvalidArgument,
	client.ErrUserNotFound:        codes.NotFound,
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
	if _, ok := status.FromError(err); ok {
		return err
	}

	for currentErr, resCode := range allErr {
		if errors.Is(err, currentErr) {
			return status.Error(resCode, err.Error())
		}
	}

	return status.Error(codes.Internal, ErrInternalServer)
}
