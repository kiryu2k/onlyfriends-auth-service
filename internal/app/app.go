package app

import (
	"context"
	"net"

	"github.com/kiryu2k/onlyfriends-auth-service/config"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/repository"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/auth"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/hasher"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/token"
	grpc2 "github.com/kiryu2k/onlyfriends-auth-service/internal/transport/grpc"
	proto "github.com/kiryu2k/onlyfriends-protos/auth"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg *config.Config) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return errors.WithMessage(err, "new logger")
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(errors.WithMessage(err, "sync logger"))
		}
	}()
	sugar := logger.Sugar()

	db, err := connectToDatabase(ctx, cfg.Database)
	if err != nil {
		return errors.WithMessage(err, "connect to db")
	}
	defer func() {
		if err := db.Close(); err != nil {
			sugar.Error(errors.WithMessage(err, "close db connection"))
		}
	}()

	authSvc := auth.New(
		hasher.New(),
		token.New(cfg.TokenSigningKey),
		repository.NewAuthUser(db),
	)

	lis, err := net.Listen("tcp", cfg.Server.Address())
	if err != nil {
		return errors.WithMessagef(err, "listen tcp %d port", cfg.Server.Port)
	}

	authTransport := grpc2.NewAuth(authSvc)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc2.ValidationMiddleware,
			grpc2.LoggingMiddleware(logger),
		),
	)
	proto.RegisterAuthServer(server, authTransport)
	healthcheck(server)

	if err := server.Serve(lis); err != nil {
		return errors.WithMessage(err, "serve grpc server")
	}

	return nil
}
