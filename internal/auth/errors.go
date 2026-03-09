package auth

import "errors"

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionInactive      = errors.New("session is no longer active")
	ErrBrowserMismatch      = errors.New("session was created from a different browser")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrPasswordMismatch     = errors.New("passwords do not match")
	ErrSamePassword         = errors.New("new password must be different from current")
	ErrInvalidCurrentPassword = errors.New("current password is incorrect")
	ErrEmailAlreadyExists   = errors.New("email already registered")
	ErrRoleNotFound         = errors.New("role not found")
	ErrDepartmentNotFound   = errors.New("department not found")
	ErrPositionNotFound     = errors.New("position not found")
	ErrUserNotFound         = errors.New("user not found")
)