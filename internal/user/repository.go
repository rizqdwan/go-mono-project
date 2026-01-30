package user

import "context"

type UserRepository interface {
	GetDetailsUser(ctx context.Context, id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	CreateUser(ctx context.Context,user *User) error
	UpdateUser(user *User) error
	DeleteUser(ctx context.Context, userId string) (err error)
}