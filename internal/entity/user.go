package entity

import (
	"time"
)

type AuthUser struct {
	Id             string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
