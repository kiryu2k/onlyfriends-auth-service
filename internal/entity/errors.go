package entity

import (
	"github.com/pkg/errors"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserNotFound      = errors.New("user not found")
)
