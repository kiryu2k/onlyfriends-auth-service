package grpc

import (
	"context"

	"github.com/kiryu2k/onlyfriends-auth-service/internal/entity"
	proto "github.com/kiryu2k/onlyfriends-protos/auth"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService interface {
	SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error)
	SignIn(ctx context.Context, req *proto.SignInRequest) (*proto.SignInResponse, error)
}

type auth struct {
	proto.UnimplementedAuthServer
	svc authService
}

func NewAuth(svc authService) auth {
	return auth{svc: svc}
}

func (t auth) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	resp, err := t.svc.SignUp(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, "sign up")
	}
	return resp, nil
}

func (t auth) SignIn(ctx context.Context, req *proto.SignInRequest) (*proto.SignInResponse, error) {
	resp, err := t.svc.SignIn(ctx, req)
	switch {
	case errors.Is(err, entity.ErrIncorrectPassword):
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, entity.ErrUserNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		return nil, errors.WithMessage(err, "sign in")
	default:
		return resp, nil
	}
}
