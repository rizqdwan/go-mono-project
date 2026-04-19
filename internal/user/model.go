package user

import (
	"time"
)

type User struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	PasswordHash string    `db:"password_hash"`
	RoleID       int64     `db:"role_id"`
	DepartmentID int64     `db:"department_id"`
	PositionID   string    `db:"position_id"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserResponse struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserFilter struct {
	Name  string
	Email string
	Role  string
}
