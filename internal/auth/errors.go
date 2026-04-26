package auth

import (
	"errors"
)

var (
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrPasswordMismatch       = errors.New("passwords do not match")
	ErrSamePassword           = errors.New("new password must differ from current password")
	ErrIncorrectPassword      = errors.New("current password is incorrect")
	ErrAccountDisabled        = errors.New("account is disabled")
	ErrSessionConflict        = errors.New("already signed in on another device, please sign out first")
	ErrSessionNotFound        = errors.New("session not found")
	ErrSessionInactive        = errors.New("session is no longer active, please login again")
	ErrBrowserMismatch        = errors.New("session was created from a different browser")
	ErrRefreshTokenExpired    = errors.New("refresh token has expired, please login again")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token, please login again")
	ErrUserInactive           = errors.New("account is inactive")
)
