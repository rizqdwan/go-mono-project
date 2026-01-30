package user

import (
	"time"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}


type Repository interface {
	CreateUser(user User) (User, error)
	GetUserByID(id string) (User, error)
	GetUserByEmail(email string) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id string) error
}