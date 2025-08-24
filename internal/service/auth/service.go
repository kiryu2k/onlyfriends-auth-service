package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/entity"
	proto "github.com/kiryu2k/onlyfriends-protos/auth"
	"github.com/pkg/errors"
)

type authRepo interface {
	InsertAuthUser(ctx context.Context, req entity.AuthUser) error
	GetAuthUserByEmail(ctx context.Context, email string) (*entity.AuthUser, error)
}

type hasher interface {
	Hash(ctx context.Context, v string) (string, error)
	VerifyHash(ctx context.Context, v string, hash string) (bool, error)
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

	err = s.repo.InsertAuthUser(ctx, entity.AuthUser{
		Id:             userId,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
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

func (s service) SignIn(ctx context.Context, req *proto.SignInRequest) (*proto.SignInResponse, error) {
	u, err := s.repo.GetAuthUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.WithMessage(err, "get auth user by email")
	}

	isEqual, err := s.hasher.VerifyHash(ctx, req.Password, u.HashedPassword)
	if err != nil {
		return nil, errors.WithMessage(err, "verify hash")
	}

	if !isEqual {
		return nil, entity.ErrIncorrectPassword
	}

	tokens, err := s.tokenSvc.GenerateTokens(ctx, entity.GenerateTokensPayload{UserId: u.Id})
	if err != nil {
		return nil, errors.WithMessage(err, "generate tokens")
	}

	return &proto.SignInResponse{
		UserId:       u.Id,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
