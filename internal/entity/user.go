package entity

type CreateAuthUserRequest struct {
	UserId         string
	Email          string
	HashedPassword string
}
