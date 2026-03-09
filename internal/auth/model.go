package auth

import "time"

type Authentication struct {
	ID 						int64 `db:"id"`
	UserID        int64 `db:"user_id"`
	RefreshToken 	string `db:"refresh_token"`
	IsActive      bool	`db:"is_active"`
	BrowserInfo   string	`db:"browser_info"`
	CreatedAt     time.Time	`db:"created_at"`
	LastActivity  time.Time	`db:"last_activity"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, min=8"`
}

type RegisterRequest struct {
	Name            string `json:"name"             validate:"required"`
	Email           string `json:"email"            validate:"required,email"`
	Password        string `json:"password"         validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	Role            string `json:"role"             validate:"required"`
	Department      string `json:"department"       validate:"required"`
	Position        string `json:"position"         validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"     validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type ChangePasswordResponse struct {
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}