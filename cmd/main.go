package main

import (
	"context"
	"log"
	"net"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/kiryu2k/onlyfriends-auth-service/config"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/repository"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/auth"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/hasher"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/token"
	grpc2 "github.com/kiryu2k/onlyfriends-auth-service/internal/transport/grpc"
	proto "github.com/kiryu2k/onlyfriends-protos/auth"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func main() {
	validate := validator.New()
	cfg, err := config.Load(validate)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "load config"))
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.Database.Dsn())
	if err != nil {
		log.Fatal(errors.WithMessage(err, "pgx connect"))
	}
	defer func() {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal(errors.WithMessage(err, "close pgx connection"))
		}
	}()

	authSvc := auth.New(
		hasher.New(),
		token.New(cfg.TokenSigningKey),
		repository.NewAuthUser(conn),
	)

	lis, err := net.Listen("tcp", cfg.Server.Address())
	if err != nil {
		log.Fatal(errors.WithMessagef(err, "listen tcp %d port", cfg.Server.Port))
	}

	authTransport := grpc2.NewAuth(authSvc)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc2.ValidationMiddleware),
	)
	proto.RegisterAuthServer(server, authTransport)

	err = server.Serve(lis)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "serve grpc server"))
	}
}
