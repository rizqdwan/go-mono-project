package user

import (
	"time"

	"github.com/rizqdwan/go-mono-project/internal/organization/department"
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

// @name UserResponse
type UserResponse struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}

// @name UserDetailsResponse
type UserDetailsResponse struct {
	ID         int64                     `json:"id"`
	Email      string                    `json:"email"`
	Name       string                    `json:"name"`
	Role       string                    `json:"role"`
	Department department.DepartmentInfo `json:"department"`
	Position   string                    `json:"position"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
}

// @name UserListResponse
type UserListResponse struct {
	ID         int64                     `json:"id"`
	Email      string                    `json:"email"`
	Name       string                    `json:"name"`
	Role       string                    `json:"role"`
	Department department.DepartmentInfo `json:"department"`
	Position   string                    `json:"position"`
	CreatedAt  time.Time                 `json:"created_at"`
}

type UserFilter struct {
	Name  string
	Email string
	Role  string
}

// @name UpdateUserDetails
type UpdateUserDetailsRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required"`
	Position string `json:"position" validate:"required"`
}

// @name ChangePasswordRequest
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"     validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// @name ChangePasswordResponse
type ChangePasswordResponse struct {
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}

// @name ResetPasswordRequest
type ResetPasswordRequest struct {
	NewPassword     string `json:"new_password"     validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// @name ResetPasswordResponse
type ResetPasswordResponse struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}
