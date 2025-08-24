package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/kiryu2k/onlyfriends-auth-service/internal/entity"
	"github.com/pkg/errors"
)

type authUser struct {
	db *pgx.Conn
}

func NewAuthUser(db *pgx.Conn) authUser {
	return authUser{db: db}
}

func (r authUser) InsertAuthUser(ctx context.Context, req entity.AuthUser) error {
	query := `
		INSERT INTO auth_user (id, email, hashed_password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, req.Id, req.Email, req.HashedPassword, req.CreatedAt, req.UpdatedAt)
	if err != nil {
		return errors.WithMessagef(err, "exec insert query: '%s'", query)
	}
	return nil
}

func (r authUser) GetAuthUserByEmail(ctx context.Context, email string) (*entity.AuthUser, error) {
	query := `
		SELECT (id, email, hashed_password, created_at, updated_at)
		FROM auth_user
		WHERE email = $1`
	u := new(entity.AuthUser)
	err := r.db.QueryRow(ctx, query, email).Scan(u)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, errors.WithMessagef(entity.ErrUserNotFound, "select query row: '%s'", query)
	case err != nil:
		return nil, errors.WithMessagef(err, "select query row: '%s'", query)
	default:
		return u, nil
	}
}
