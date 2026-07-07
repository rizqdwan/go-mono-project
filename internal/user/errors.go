package user

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrEmailAlreadyExists       = errors.New("email already registered")
	ErrRoleNotFound             = errors.New("role not found")
	ErrDepartmentNotFound       = errors.New("department not found")
	ErrPositionNotFound         = errors.New("position not found")
	ErrIncorrectPassword        = errors.New("current password is incorrect")
	ErrCannotDeleteUser         = errors.New("cannot delete user")
	ErrCannotDeleteSelf         = errors.New("cannot delete self")
	ErrCannotDeleteAdmin        = errors.New("cannot delete admin")
	ErrDifferentDepartment      = errors.New("different department")
	ErrUserHasActiveProjects    = errors.New("user has active projects")
	ErrCannotAssignProjectAdmin = errors.New("cannot assign project admin role")
	ErrCannotManageSelf         = errors.New("cannot manage self")
)
