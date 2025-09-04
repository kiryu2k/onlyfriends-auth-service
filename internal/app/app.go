package app

import (
	"context"
	"database/sql"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kiryu2k/onlyfriends-auth-service/config"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/repository"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/auth"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/hasher"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/service/token"
	grpc2 "github.com/kiryu2k/onlyfriends-auth-service/internal/transport/grpc"
	"github.com/kiryu2k/onlyfriends-auth-service/migrations"
	proto "github.com/kiryu2k/onlyfriends-protos/auth"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg *config.Config) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return errors.WithMessage(err, "new logger")
	}
	//defer func() {
	//	if err := logger.Sync(); err != nil {
	//		panic(errors.WithMessage(err, "sync logger"))
	//	}
	//}()
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
		grpc.UnaryInterceptor(grpc2.ValidationMiddleware),
	)
	proto.RegisterAuthServer(server, authTransport)

	if err := server.Serve(lis); err != nil {
		return errors.WithMessage(err, "serve grpc server")
	}

	return nil
}

func connectToDatabase(ctx context.Context, cfg config.Database) (*sqlx.DB, error) {
	connCfg, err := pgx.ParseConfig(cfg.Dsn())
	if err != nil {
		return nil, errors.WithMessage(err, "parse config")
	}

	sqlDb := stdlib.OpenDB(*connCfg)
	pgDb := sqlx.NewDb(sqlDb, "pgx")
	if err := pgDb.PingContext(ctx); err != nil {
		return nil, errors.WithMessage(err, "ping db")
	}

	if _, err := pgDb.ExecContext(ctx, `CREATE SCHEMA IF NOT EXISTS `+cfg.Schema); err != nil {
		return nil, errors.WithMessage(err, "exec create schema query")
	}

	if err := migrate(ctx, pgDb.DB); err != nil {
		return nil, errors.WithMessage(err, "migrate db")
	}

	return pgDb, nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(migrations.Files)
	if err := goose.SetDialect("postgres"); err != nil {
		return errors.WithMessage(err, "goose set dialect")
	}
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return errors.WithMessage(err, "goose up migrations")
	}
	return nil
}
