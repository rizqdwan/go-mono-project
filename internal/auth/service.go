package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/rizqdwan/go-mono-project/infrastructure/security"
	"github.com/rizqdwan/go-mono-project/internal/user"
	"github.com/rizqdwan/go-mono-project/pkg/token"
)

type Service interface {
	Login(ctx context.Context, req LoginRequest, browserInfo string) (*TokenResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*user.UserResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	RenewToken(ctx context.Context, req RefreshTokenRequest, browserInfo string) (*TokenResponse, error)
	ChangePassword(ctx context.Context, userID int64, req ChangePasswordRequest) (ChangePasswordResponse, error)
}

const maxActiveRefreshTokens = 1

type service struct {
	authRepo    Repository
	userRepo    user.Repository
	tokenSvc    *token.Service
	passwordSvc security.PasswordHash
}

func NewService(authRepo Repository, userRepo user.Repository, tokenSvc *token.Service, passwordSvc security.PasswordHash) Service {
	s := &service{
		authRepo:    authRepo,
		userRepo:    userRepo,
		tokenSvc:    tokenSvc,
		passwordSvc: passwordSvc,
	}

	go s.startSessionCleanup()
	return s
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*user.UserResponse, error) {
	if req.Password != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	u, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && u != nil {
		return nil, ErrEmailAlreadyRegistered
	}
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	roleID, err := s.userRepo.FindRoleByName(ctx, req.Role)
	if err != nil {
		if errors.Is(err, user.ErrRoleNotFound) {
			return nil, user.ErrRoleNotFound
		}
		return nil, err
	}

	deptID, err := s.userRepo.FindDepartmentByLabel(ctx, req.Department)
	if err != nil {
		return nil, user.ErrDepartmentNotFound
	}

	posID, err := s.userRepo.FindPositionByLabel(ctx, req.Position)
	if err != nil {
		return nil, user.ErrPositionNotFound
	}

	hashed, err := s.passwordSvc.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	u = &user.User{
		Name:         req.Name,
		Email:        email,
		PasswordHash: hashed,
		RoleID:       roleID,
		DepartmentID: deptID,
		PositionID:   posID,
		IsActive:     true,
	}

	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	return &user.UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Role:       req.Role,
		Department: req.Department,
		Position:   posID,
		CreatedAt:  u.CreatedAt,
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest, browserInfo string) (*TokenResponse, error) {
	u, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !u.IsActive {
		return nil, ErrUserInactive
	}

	match, err := s.passwordSvc.Compare(req.Password, u.PasswordHash)
	if err != nil || !match {
		return nil, ErrInvalidCredentials
	}

	activeSessions, err := s.authRepo.FindActiveSessionsByUserID(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	for _, session := range activeSessions {
		if session.BrowserInfo != browserInfo {
			return nil, ErrSessionConflict
		}
	}

	roleName, err := s.userRepo.FindRoleNameByID(ctx, u.RoleID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(u.ID, u.Email, roleName)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}

	if err := s.enforceMaxActiveSessions(ctx, u.ID); err != nil {
		return nil, err
	}

	auth, findErr := s.authRepo.FindSessionByUserAndBrowser(ctx, u.ID, browserInfo)
	if findErr != nil && !errors.Is(findErr, ErrSessionNotFound) {
		return nil, findErr
	}

	if auth == nil {
		newSession := &Authentication{
			UserID:       u.ID,
			RefreshToken: refreshToken,
			IsActive:     true,
			BrowserInfo:  browserInfo,
		}
		if err := s.authRepo.CreateSession(ctx, newSession); err != nil {
			return nil, err
		}
	} else {
		auth.RefreshToken = refreshToken
		auth.IsActive = true
		if err := s.authRepo.UpdateSession(ctx, auth); err != nil {
			return nil, err
		}
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokenSvc.AccessTTL().Seconds()),
	}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	auth, err := s.authRepo.FindSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return ErrSessionNotFound
	}

	if err := s.authRepo.DeactivateSession(ctx, auth.ID); err != nil {
		return ErrSessionNotFound
	}

	return nil
}

func (s *service) RenewToken(ctx context.Context, req RefreshTokenRequest, browserInfo string) (*TokenResponse, error) {
	auth, err := s.authRepo.FindSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if !auth.IsActive {
		return nil, ErrSessionInactive
	}

	if auth.BrowserInfo != browserInfo {
		_ = s.authRepo.DeactivateSession(ctx, auth.ID)
		return nil, ErrBrowserMismatch
	}

	claims, err := s.tokenSvc.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		_ = s.authRepo.DeactivateSession(ctx, auth.ID)
		switch {
		case errors.Is(err, token.ErrExpiredToken):
			return nil, ErrRefreshTokenExpired
		default:
			return nil, ErrInvalidRefreshToken
		}
	}

	u, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}
	if !u.IsActive {
		return nil, ErrUserInactive
	}

	roleName, err := s.userRepo.FindRoleNameByID(ctx, u.RoleID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(u.ID, u.Email, roleName)
	if err != nil {
		return nil, err
	}

	refreshToken := req.RefreshToken
	if time.Until(claims.ExpiresAt.Time) < 24*time.Hour {
		refreshToken, err = s.tokenSvc.GenerateRefreshToken(u.ID)
		if err != nil {
			return nil, err
		}
		auth.RefreshToken = refreshToken
	}

	if err := s.enforceMaxActiveSessions(ctx, u.ID); err != nil {
		return nil, err
	}

	if err := s.authRepo.UpdateSession(ctx, auth); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokenSvc.AccessTTL().Seconds()),
	}, nil
}

func (s *service) ChangePassword(ctx context.Context, userID int64, req ChangePasswordRequest) (ChangePasswordResponse, error) {
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ChangePasswordResponse{}, err
	}

	match, err := s.passwordSvc.Compare(req.OldPassword, existingUser.PasswordHash)
	if err != nil || !match {
		return ChangePasswordResponse{}, ErrInvalidCredentials
	}

	newHash, err := s.passwordSvc.Hash(req.NewPassword)
	if err != nil {
		return ChangePasswordResponse{}, err
	}
	updatedAt, err := s.userRepo.UpdatePassword(ctx, userID, newHash)
	if err != nil {
		return ChangePasswordResponse{}, err
	}

	return ChangePasswordResponse{
		Email:     existingUser.Email,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *service) enforceMaxActiveSessions(ctx context.Context, userID int64) error {
	// fix: consistent method name
	sessions, err := s.authRepo.FindActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if len(sessions) < maxActiveRefreshTokens {
		return nil
	}

	excess := len(sessions) - maxActiveRefreshTokens + 1
	for _, session := range sessions[:excess] {
		if err := s.authRepo.DeactivateSession(ctx, session.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) startSessionCleanup() {
	interval := s.tokenSvc.RefreshTTL()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		cutoff := time.Now().Add(-interval)

		sessions, err := s.authRepo.FindExpiredActiveSessions(ctx, cutoff, 100)
		if err != nil {
			continue
		}

		for _, session := range sessions {
			if _, err := s.tokenSvc.ValidateRefreshToken(session.RefreshToken); err != nil {
				_ = s.authRepo.DeactivateSession(ctx, session.ID)
			}
		}

		_ = s.authRepo.DeleteInactiveSessionsOlderThan(ctx, time.Now().AddDate(0, 0, -7))
	}
}
