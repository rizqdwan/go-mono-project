package auth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/pkg/response"
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

func mapError(c *echo.Context, err error) error {
	switch err {
	case ErrInvalidCredentials, ErrIncorrectPassword:
		return response.Error(c, http.StatusUnauthorized, err.Error(), nil)
	case ErrEmailAlreadyRegistered, ErrSessionConflict:
		return response.Error(c, http.StatusConflict, err.Error(), nil)
	case ErrPasswordMismatch, ErrSamePassword:
		return response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
	case ErrSessionNotFound, ErrSessionInactive,
		ErrRefreshTokenExpired, ErrInvalidRefreshToken,
		ErrBrowserMismatch:
		return response.Error(c, http.StatusUnauthorized, err.Error(), nil)
	case ErrUserInactive, ErrAccountDisabled:
		return response.Error(c, http.StatusForbidden, err.Error(), nil)
	default:
		return response.Error(c, http.StatusInternalServerError, "internal server error", nil)
	}
}
