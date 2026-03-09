package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrRoleNotFound       = errors.New("role not found")
	ErrDepartmentNotFound = errors.New("department not found")
	ErrPositionNotFound   = errors.New("position not found")
)