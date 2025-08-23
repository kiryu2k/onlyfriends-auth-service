package repository

import (
	"context"

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

func (r authUser) InsertAuthUser(ctx context.Context, req entity.CreateAuthUserRequest) error {
	query := `INSERT INTO auth_user (id, email, hashed_password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, req.UserId, req.Email, req.HashedPassword)
	if err != nil {
		return errors.WithMessagef(err, "exec insert query: '%s'", query)
	}
	return nil
}
