package app

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kiryu2k/onlyfriends-auth-service/config"
	"github.com/kiryu2k/onlyfriends-auth-service/migrations"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

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
