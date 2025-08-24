package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/entity"
	proto "github.com/kiryu2k/onlyfriends-protos/auth"
	"github.com/pkg/errors"
)

type authRepo interface {
	InsertAuthUser(ctx context.Context, req entity.CreateAuthUserRequest) error
}

type hasher interface {
	Hash(ctx context.Context, v string) (string, error)
}

type tokenService interface {
	GenerateTokens(ctx context.Context, payload entity.GenerateTokensPayload) (*entity.GenerateTokensResult, error)
}

type service struct {
	hasher   hasher
	tokenSvc tokenService
	repo     authRepo
}

func New(hasher hasher, tokenSvc tokenService, repo authRepo) service {
	return service{
		hasher:   hasher,
		tokenSvc: tokenSvc,
		repo:     repo,
	}
}

func (s service) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	hashedPassword, err := s.hasher.Hash(ctx, req.Password)
	if err != nil {
		return nil, errors.WithMessage(err, "hash")
	}
	userId := uuid.NewString()

	err = s.repo.InsertAuthUser(ctx, entity.CreateAuthUserRequest{
		UserId:         userId,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "insert auth user")
	}

	tokens, err := s.tokenSvc.GenerateTokens(ctx, entity.GenerateTokensPayload{UserId: userId})
	if err != nil {
		return nil, errors.WithMessage(err, "generate tokens")
	}

	return &proto.SignUpResponse{
		UserId:       userId,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
