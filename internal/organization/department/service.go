package department

import (
	"context"
	"errors"

	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
	"github.com/rizqdwan/go-mono-project/internal/organization/group"
	"github.com/rizqdwan/go-mono-project/pkg/pagination"
)

type Service interface {
	ListDepartments(ctx context.Context, p pagination.Params) ([]DepartmentResponse, int, error)
	CreateDepartment(ctx context.Context, req NewDepartmentRequest) (DepartmentResponse, error)
	UpdateDepartment(ctx context.Context, id int64, req UpdateDepartmentRequest) (DepartmentResponse, error)
	DeleteDepartment(ctx context.Context, id int64) error
}

type service struct {
	departmentRepo Repository
	groupRepo      group.Repository
}

func NewService(departmentRepo Repository, groupRepo group.Repository) Service {
	s := &service{
		departmentRepo: departmentRepo,
		groupRepo:      groupRepo,
	}

	return s
}

func (s *service) ListDepartments(ctx context.Context, p pagination.Params) ([]DepartmentResponse, int, error) {
	department, total, err := s.departmentRepo.FindAllDepartments(ctx, p)
	if err != nil {
		return nil, 0, commonErrors.InternalServerError("failed to fetch departments", err)
	}
	return department, total, nil
}

func (s *service) CreateDepartment(ctx context.Context, req NewDepartmentRequest) (DepartmentResponse, error) {
	exisiting, err := s.departmentRepo.FindDepartmentByLabel(ctx, req.Label)
	if err == nil && exisiting != nil {
		return DepartmentResponse{}, commonErrors.ConflictError(ErrDepartmentLabelExists.Error(), ErrDepartmentLabelExists)
	}

	if err != nil && !errors.Is(err, ErrDepartmentNotFound) {
		return DepartmentResponse{}, commonErrors.InternalServerError("failed to check department label", err)
	}

	g, err := s.groupRepo.FindGroupByID(ctx, req.GroupID)
	if err != nil {
		if errors.Is(err, group.ErrGroupNotFound) {
			return DepartmentResponse{}, commonErrors.NotFoundError(group.ErrGroupNotFound.Error(), group.ErrGroupNotFound)
		}
		return DepartmentResponse{}, commonErrors.InternalServerError("failed to fetch group", err)
	}

	d := &Department{
		Label:   req.Label,
		Name:    req.Name,
		GroupID: req.GroupID,
	}

	if err := s.departmentRepo.CreateDepartment(ctx, d); err != nil {
		return DepartmentResponse{}, commonErrors.InternalServerError("failed to create department", err)
	}

	return DepartmentResponse{
		ID:    d.ID,
		Label: d.Label,
		Name:  d.Name,
		Group: group.GroupInfo{
			Label: g.Label,
			Name:  g.Name,
		},
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}, nil
}

func (s *service) UpdateDepartment(ctx context.Context, id int64, req UpdateDepartmentRequest) (DepartmentResponse, error) {
	dept, err := s.departmentRepo.FindDepartmentByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDepartmentNotFound) {
			return DepartmentResponse{}, commonErrors.NotFoundError(ErrDepartmentNotFound.Error(), ErrDepartmentNotFound)
		}
		return DepartmentResponse{}, commonErrors.InternalServerError("failed to fetch department", err)
	}

	if req.Label != dept.Label {
		existing, err := s.departmentRepo.FindDepartmentByLabel(ctx, req.Label)

		if err == nil && existing != nil {
			return DepartmentResponse{}, commonErrors.ConflictError(ErrDepartmentLabelExists.Error(), ErrDepartmentLabelExists)
		}

		if err != nil && !errors.Is(err, ErrDepartmentNotFound) {
			return DepartmentResponse{}, commonErrors.InternalServerError("failed to check department label", err)
		}
	}

	g, err := s.groupRepo.FindGroupByID(ctx, req.GroupID)
	if err != nil {
		if errors.Is(err, group.ErrGroupNotFound) {
			return DepartmentResponse{}, commonErrors.NotFoundError(group.ErrGroupNotFound.Error(), group.ErrGroupNotFound)
		}
		return DepartmentResponse{}, commonErrors.InternalServerError("failed to fetch group", err)
	}

	dept.Label = req.Label
	dept.Name = req.Name
	dept.GroupID = req.GroupID

	updatedAt, err := s.departmentRepo.UpdateDepartment(ctx, dept)
	if err != nil {
		return DepartmentResponse{}, commonErrors.InternalServerError("failed to update department", err)
	}

	return DepartmentResponse{
		ID:    dept.ID,
		Label: dept.Label,
		Name:  dept.Name,
		Group: group.GroupInfo{
			Label: g.Label,
			Name:  g.Name,
		},
		CreatedAt: dept.CreatedAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *service) DeleteDepartment(ctx context.Context, id int64) error {
	_, err := s.departmentRepo.FindDepartmentByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDepartmentNotFound) {
			return commonErrors.NotFoundError(ErrDepartmentNotFound.Error(), ErrDepartmentNotFound)
		}
		return commonErrors.InternalServerError("failed to fetch department", err)
	}

	hasUsers, err := s.departmentRepo.HasActiveUsers(ctx, id)
	if err != nil {
		return commonErrors.InternalServerError("failed to check department users", err)
	}
	if hasUsers {
		return commonErrors.ConflictError(ErrDepartmentHasActiveUsers.Error(), ErrDepartmentHasActiveUsers)
	}

	if err := s.departmentRepo.DeleteDepartment(ctx, id); err != nil {
		return commonErrors.InternalServerError("failed to delete department", err)
	}

	return nil
}
