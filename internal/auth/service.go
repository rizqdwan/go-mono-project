package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/rizqdwan/go-mono-project/infrastructure/security"
	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
	"github.com/rizqdwan/go-mono-project/internal/organization/department"
	"github.com/rizqdwan/go-mono-project/internal/user"
	"github.com/rizqdwan/go-mono-project/pkg/token"
)

type Service interface {
	Login(ctx context.Context, req LoginRequest, browserInfo string) (*TokenResponse, error)
	Register(ctx context.Context, adminID int64, req RegisterRequest) (*RegisterResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	RenewToken(ctx context.Context, req RefreshTokenRequest, browserInfo string) (*TokenResponse, error)
}

const maxActiveRefreshTokens = 1

type service struct {
	authRepo    Repository
	userRepo    user.Repository
	tokenSvc    *token.Service
	passwordSvc security.PasswordHash
}

func NewService(authRepo Repository, userRepo user.Repository, tokenSvc *token.Service, passwordSvc security.PasswordHash) Service {
	return &service{
		authRepo:    authRepo,
		userRepo:    userRepo,
		tokenSvc:    tokenSvc,
		passwordSvc: passwordSvc,
	}
}

func (s *service) Register(ctx context.Context, adminID int64, req RegisterRequest) (*RegisterResponse, error) {
	if req.Password != req.ConfirmPassword {
		return nil, commonErrors.InvariantError(ErrPasswordMismatch.Error(), ErrPasswordMismatch)
	}

	admin, err := s.userRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch admin", err)
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	_, err = s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, commonErrors.ConflictError(ErrEmailAlreadyRegistered.Error(), ErrEmailAlreadyRegistered)
	}
	if !errors.Is(err, user.ErrUserNotFound) {
		return nil, commonErrors.InternalServerError("failed to check existing email", err)
	}

	roleID, err := s.userRepo.FindRoleByName(ctx, req.Role)
	if err != nil {
		if errors.Is(err, user.ErrRoleNotFound) {
			return nil, commonErrors.NotFoundError(user.ErrRoleNotFound.Error(), err)
		}
		return nil, commonErrors.InternalServerError("failed to find role", err)
	}

	posID, err := s.userRepo.FindPositionByID(ctx, req.Position)
	if err != nil {
		if errors.Is(err, user.ErrPositionNotFound) {
			return nil, commonErrors.NotFoundError(user.ErrPositionNotFound.Error(), user.ErrPositionNotFound)
		}
		return nil, commonErrors.InternalServerError("failed to fetch position", err)
	}

	hashed, err := s.passwordSvc.Hash(req.Password)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to hash password", err)
	}

	u := &user.User{
		Name:         req.Name,
		Email:        email,
		PasswordHash: hashed,
		RoleID:       roleID,
		DepartmentID: admin.DepartmentID,
		PositionID:   posID,
		IsActive:     true,
	}

	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return nil, commonErrors.InternalServerError("failed to create user", err)
	}

	details, err := s.userRepo.FindUserDetailsByID(ctx, u.ID)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch registered user", err)
	}

	return &RegisterResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
		Role:  req.Role,
		Department: department.DepartmentInfo{
			Label: details.Department.Label,
			Name:  details.Department.Name,
		},
		Position:  posID,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest, browserInfo string) (*TokenResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, commonErrors.UnauthorizedError(ErrInvalidCredentials.Error(), ErrInvalidCredentials)
		}
		return nil, commonErrors.InternalServerError("failed to fetch user", err)
	}

	if !u.IsActive {
		return nil, commonErrors.ForbiddenError(ErrUserInactive.Error(), ErrUserInactive)
	}

	match, err := s.passwordSvc.Compare(req.Password, u.PasswordHash)
	if err != nil || !match {
		return nil, commonErrors.UnauthorizedError(ErrInvalidCredentials.Error(), ErrInvalidCredentials)
	}

	activeSessions, err := s.authRepo.FindActiveSessionsByUserID(ctx, u.ID)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch active sessions", err)
	}

	for _, session := range activeSessions {
		if session.BrowserInfo != browserInfo {
			return nil, commonErrors.ConflictError(ErrSessionConflict.Error(), ErrSessionConflict)
		}
	}

	roleName, err := s.userRepo.FindRoleNameByID(ctx, u.RoleID)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch role name", err)
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(u.ID, u.Email, roleName)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to generate access token", err)
	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to generate refresh token", err)
	}

	existingSession, findErr := s.authRepo.FindSessionByUserAndBrowser(ctx, u.ID, browserInfo)
	if findErr != nil && !errors.Is(findErr, ErrSessionNotFound) {
		return nil, commonErrors.InternalServerError("failed to find session", findErr)
	}

	expiresAt := time.Now().Add(s.tokenSvc.RefreshTTL())

	if existingSession == nil {
		if err := s.enforceMaxActiveSessions(ctx, u.ID); err != nil {
			return nil, err
		}
		newSession := &Authentication{
			UserID:       u.ID,
			RefreshToken: refreshToken,
			IsActive:     true,
			BrowserInfo:  browserInfo,
			ExpiresAt:    expiresAt,
		}
		if err := s.authRepo.CreateSession(ctx, newSession); err != nil {
			return nil, commonErrors.InternalServerError("failed to create session", err)
		}
	} else {
		existingSession.RefreshToken = refreshToken
		existingSession.ExpiresAt = expiresAt
		if err := s.authRepo.UpdateSession(ctx, existingSession); err != nil {
			return nil, commonErrors.InternalServerError("failed to update session", err)
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
		return commonErrors.UnauthorizedError(ErrSessionNotFound.Error(), ErrSessionNotFound)
	}

	if err := s.authRepo.DeactivateSession(ctx, auth.ID); err != nil {
		return commonErrors.InternalServerError("failed to deactivate session", err)
	}

	return nil
}

func (s *service) RenewToken(ctx context.Context, req RefreshTokenRequest, browserInfo string) (*TokenResponse, error) {
	claims, err := s.tokenSvc.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, commonErrors.UnauthorizedError(ErrInvalidRefreshToken.Error(), ErrInvalidRefreshToken)
	}

	auth, err := s.authRepo.FindSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, commonErrors.UnauthorizedError(ErrSessionNotFound.Error(), ErrSessionNotFound)
		}
		return nil, commonErrors.InternalServerError("failed to fetch session", err)
	}

	if auth.BrowserInfo != browserInfo {
		_ = s.authRepo.DeactivateSession(ctx, auth.ID)
		return nil, commonErrors.UnauthorizedError(ErrBrowserMismatch.Error(), ErrBrowserMismatch)
	}

	u, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch user", err)
	}

	if !u.IsActive {
		return nil, commonErrors.ForbiddenError(ErrUserInactive.Error(), ErrUserInactive)
	}

	roleName, err := s.userRepo.FindRoleNameByID(ctx, u.RoleID)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch role name", err)
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(
		u.ID,
		u.Email,
		roleName,
	)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to generate access token", err)
	}

	refreshToken := req.RefreshToken

	if time.Until(claims.ExpiresAt.Time) <= 24*time.Hour {
		expiresAt := time.Now().Add(s.tokenSvc.RefreshTTL())

		refreshToken, err = s.tokenSvc.GenerateRefreshToken(u.ID)
		if err != nil {
			return nil, commonErrors.InternalServerError("failed to generate refresh token", err)
		}
		auth.RefreshToken = refreshToken
		auth.ExpiresAt = expiresAt
	}

	if err := s.authRepo.UpdateSession(ctx, auth); err != nil {
		return nil, commonErrors.InternalServerError("failed to update session", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokenSvc.AccessTTL().Seconds()),
	}, nil
}

func (s *service) enforceMaxActiveSessions(ctx context.Context, userID int64) error {
	sessions, err := s.authRepo.FindActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return commonErrors.InternalServerError("failed to fetch active sessions", err)
	}

	if len(sessions) < maxActiveRefreshTokens {
		return nil
	}

	excess := len(sessions) - maxActiveRefreshTokens + 1
	for _, session := range sessions[:excess] {
		if err := s.authRepo.DeactivateSession(ctx, session.ID); err != nil {
			return commonErrors.InternalServerError("failed to deactivate session", err)
		}
	}
	return nil
}
