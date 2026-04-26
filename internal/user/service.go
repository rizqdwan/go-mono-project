package user

import (
	"context"

	"github.com/rizqdwan/go-mono-project/infrastructure/security"
	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
	"github.com/rizqdwan/go-mono-project/pkg/token"
)

type SessionDeactivator interface {
	DeactivateAllUserSessions(ctx context.Context, userID int64) error
}

type Service interface {
	UserDetails(ctx context.Context, userID int64) (*UserDetailsResponse, error)
	ChangePassword(ctx context.Context, userID int64, req ChangePasswordRequest) (ChangePasswordResponse, error)
	ResetPassword(ctx context.Context, targetUserID int64, req ResetPasswordRequest) (ResetPasswordResponse, error)
}

type service struct {
	userRepo    Repository
	tokenSvc    *token.Service
	passwordSvc security.PasswordHash
	sessionSvc  SessionDeactivator
}

func NewService(userRepo Repository, tokenSvc *token.Service, passwordSvc security.PasswordHash, sessionSvc SessionDeactivator) Service {
	s := &service{
		userRepo:    userRepo,
		tokenSvc:    tokenSvc,
		passwordSvc: passwordSvc,
		sessionSvc:  sessionSvc,
	}

	return s
}

func (s *service) UserDetails(ctx context.Context, userID int64) (*UserDetailsResponse, error) {
	user, err := s.userRepo.FindUserDetailsByID(ctx, userID)
	if err != nil {
		return nil, commonErrors.NotFoundError(ErrUserNotFound.Error(), err)
	}

	return user, nil
}

func (s *service) ChangePassword(ctx context.Context, userID int64, req ChangePasswordRequest) (ChangePasswordResponse, error) {
	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ChangePasswordResponse{}, commonErrors.NotFoundError(ErrUserNotFound.Error(), err)
	}

	match, err := s.passwordSvc.Compare(req.OldPassword, existingUser.PasswordHash)
	if err != nil || !match {
		return ChangePasswordResponse{}, commonErrors.UnauthorizedError(ErrIncorrectPassword.Error(), ErrIncorrectPassword)
	}

	newHash, err := s.passwordSvc.Hash(req.NewPassword)
	if err != nil {
		return ChangePasswordResponse{}, commonErrors.InternalServerError("failed to hash password", err)
	}

	updatedAt, err := s.userRepo.UpdatePassword(ctx, userID, newHash)
	if err != nil {
		return ChangePasswordResponse{}, commonErrors.InternalServerError("failed to update password", err)
	}

	return ChangePasswordResponse{
		Email:     existingUser.Email,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *service) ResetPassword(ctx context.Context, targetUserID int64, req ResetPasswordRequest) (ResetPasswordResponse, error) {
	targetUser, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return ResetPasswordResponse{}, commonErrors.NotFoundError(ErrUserNotFound.Error(), err)
	}

	newHsh, err := s.passwordSvc.Hash(req.NewPassword)
	if err != nil {
		return ResetPasswordResponse{}, commonErrors.InternalServerError("failed to hash password", err)
	}

	updatedAt, err := s.userRepo.UpdatePassword(ctx, targetUserID, newHsh)
	if err != nil {
		return ResetPasswordResponse{}, commonErrors.InternalServerError("failed to update password", err)
	}

	if err := s.sessionSvc.DeactivateAllUserSessions(ctx, targetUserID); err != nil {
		return ResetPasswordResponse{}, commonErrors.InternalServerError("failed to invalidate user sessions", err)
	}

	return ResetPasswordResponse{
		UserID:    targetUser.ID,
		Email:     targetUser.Email,
		UpdatedAt: updatedAt,
	}, nil
}
