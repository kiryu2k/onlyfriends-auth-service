package token

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/entity"
	"github.com/pkg/errors"
)

const (
	accessTokenTtl  = 15 * time.Minute
	refreshTokenTtl = 7 * 24 * time.Hour

	issuer = "onlyfriends-auth-service"
)

type service struct {
	signingKey []byte
}

func New(signingKey string) service {
	return service{signingKey: []byte(signingKey)}
}

type claims struct {
	entity.GenerateTokensPayload
	jwt.RegisteredClaims
}

func (s service) GenerateTokens(_ context.Context, payload entity.GenerateTokensPayload) (*entity.GenerateTokensResult, error) {
	accessToken, err := s.accessToken(payload)
	if err != nil {
		return nil, errors.WithMessage(err, "access token")
	}
	refreshToken, err := s.refreshToken(payload)
	if err != nil {
		return nil, errors.WithMessage(err, "refresh token")
	}

	return &entity.GenerateTokensResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s service) accessToken(payload entity.GenerateTokensPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		GenerateTokensPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   "",
			Audience:  []string{"onlyfriends-gate-service"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTtl)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "",
		},
	})
	result, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", errors.WithMessage(err, "token signed string")
	}
	return result, nil
}

func (s service) refreshToken(payload entity.GenerateTokensPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		GenerateTokensPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   "",
			Audience:  []string{issuer},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTtl)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "",
		},
	})
	result, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", errors.WithMessage(err, "token signed string")
	}
	return result, nil
}
