package department

import "errors"

var (
	ErrDepartmentNotFound       = errors.New("department not found")
	ErrDepartmentLabelExists    = errors.New("department label already exists")
	ErrDepartmentHasActiveUsers = errors.New("department has active users")
)
