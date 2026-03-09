package auth

import "context"

type Service interface {
	Login(ctx context.Context, req LoginRequest, browserInfo, ipAddress string) (*TokenResponse, error)
	Register(ctx context.Context, req RegisterRequest) error
	Logout(ctx context.Context, refreshToken string) error
	RenewToken(ctx context.Context, req RefreshTokenRequest, browserInfo string) (*TokenResponse, error)
	ChangePassword(ctx context.Context, userID int64, req ChangePasswordRequest) error
}