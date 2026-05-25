package group

import "errors"

var (
	ErrGroupNotFound            = errors.New("group not found")
	ErrGroupLabelExists         = errors.New("group label already exists")
	ErrGroupHasActiveDepartment = errors.New("group has active department")
)
