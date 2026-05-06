package user

import (
	"context"

	"github.com/rizqdwan/go-mono-project/config"
	"github.com/rizqdwan/go-mono-project/infrastructure/security"
	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
	"github.com/rizqdwan/go-mono-project/internal/user/role"
	"github.com/rizqdwan/go-mono-project/pkg/token"
)

type SessionDeactivator interface {
	DeactivateAllUserSessions(ctx context.Context, userID int64) error
}

type Service interface {
	ListUser(ctx context.Context, userID int64) ([]UserListResponse, error)
	UserDetails(ctx context.Context, userID int64) (*UserDetailsResponse, error)
	//UpdateUserDetails(ctx context.Context, user *UserDetailsResponse) (*UserDetailsResponse, error)
	ChangePassword(ctx context.Context, userID int64, req ChangePasswordRequest) (ChangePasswordResponse, error)
	ResetPassword(ctx context.Context, targetUserID int64, req ResetPasswordRequest) (ResetPasswordResponse, error)
	DeleteUser(ctx context.Context, adminID int64, targetUserID int64) error
}

type service struct {
	userRepo    Repository
	tokenSvc    *token.Service
	passwordSvc security.PasswordHash
	sessionSvc  SessionDeactivator
	roles       config.RoleConfig
}

func NewService(userRepo Repository, tokenSvc *token.Service, passwordSvc security.PasswordHash, sessionSvc SessionDeactivator, roles config.RoleConfig) Service {
	s := &service{
		userRepo:    userRepo,
		tokenSvc:    tokenSvc,
		passwordSvc: passwordSvc,
		sessionSvc:  sessionSvc,
		roles:       roles,
	}

	return s
}

func (s *service) ListUser(ctx context.Context, userID int64) ([]UserListResponse, error) {
	admin, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, commonErrors.NotFoundError(ErrUserNotFound.Error(), err)
	}

	users, err := s.userRepo.FindUsersByDepartmentID(ctx, admin.DepartmentID)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch users", err)
	}

	for i := range users {
		users[i].Role = role.FormatRoleName(users[i].Role)
	}

	return users, nil
}

func (s *service) UserDetails(ctx context.Context, userID int64) (*UserDetailsResponse, error) {
	u, err := s.userRepo.FindUserDetailsByID(ctx, userID)
	if err != nil {
		return nil, commonErrors.NotFoundError(ErrUserNotFound.Error(), err)
	}

	u.Role = role.FormatRoleName(u.Role)

	return u, nil
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

func (s *service) DeleteUser(ctx context.Context, adminID int64, targetUserID int64) error {
	if adminID == targetUserID {
		return commonErrors.ForbiddenError(ErrCannotDeleteSelf.Error(), ErrCannotDeleteSelf)
	}

	admin, err := s.userRepo.FindByID(ctx, adminID)
	if err != nil {
		return commonErrors.InternalServerError("failed to fetch admin user", err)
	}

	targetUser, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return commonErrors.NotFoundError(ErrUserNotFound.Error(), err)
	}

	targetRoleName, err := s.userRepo.FindRoleNameByID(ctx, targetUser.RoleID)
	if err != nil {
		return commonErrors.InternalServerError("failed to fetch role name", err)
	}

	if targetRoleName == s.roles.ProjectAdmin {
		return commonErrors.ForbiddenError(ErrCannotDeleteAdmin.Error(), ErrCannotDeleteAdmin)
	}

	if admin.DepartmentID != targetUser.DepartmentID {
		return commonErrors.ForbiddenError(ErrDifferentDepartment.Error(), ErrDifferentDepartment)
	}

	hasActive, err := s.userRepo.HasActiveProjectAssignments(ctx, targetUser.ID)
	if err != nil {
		return commonErrors.InternalServerError("failed to check project assignments", err)
	}

	if hasActive {
		return commonErrors.ConflictError(ErrUserHasActiveProjects.Error(), ErrUserHasActiveProjects)
	}

	if err := s.userRepo.DeactivateUser(ctx, targetUserID); err != nil {
		return commonErrors.InternalServerError("failed to deactivate user", err)
	}

	if err := s.sessionSvc.DeactivateAllUserSessions(ctx, targetUserID); err != nil {
		return commonErrors.InternalServerError("failed to revoke user sessions", err)
	}

	return nil
}
