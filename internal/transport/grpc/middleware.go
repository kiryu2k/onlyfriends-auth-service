package grpc

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validatable interface {
	Validate() error
}

func ValidationMiddleware(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp any, err error) {
	request, ok := req.(validatable)
	if !ok {
		return handler(ctx, req)
	}

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.WithMessage(err, "validate request body").Error())
	}

	return handler(ctx, req)
}
