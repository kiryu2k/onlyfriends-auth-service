package grpc

import (
	"context"

	proto "github.com/kiryu2k/onlyfriends-protos/auth"
	"github.com/pkg/errors"
)

type authService interface {
	SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error)
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
