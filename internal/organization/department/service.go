package department

import (
	"context"

	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
	"github.com/rizqdwan/go-mono-project/internal/organization/group"
)

type Service interface {
	ListDepartments(ctx context.Context) ([]DepartmentResponse, error)
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

func (s *service) ListDepartments(ctx context.Context) ([]DepartmentResponse, error) {
	department, err := s.departmentRepo.FindAllDepartments(ctx)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch departments", err)
	}
	return department, nil
}

func (s *service) CreateDepartment(ctx context.Context, req NewDepartmentRequest) (DepartmentResponse, error) {
	exisiting, err := s.departmentRepo.FindDepartmentByLabel(ctx, req.Label)
	if err == nil && exisiting != nil {
		return DepartmentResponse{}, commonErrors.ConflictError(ErrDepartmentLabelExists.Error(), ErrDepartmentLabelExists)
	}

	g, err := s.groupRepo.FindGroupByID(ctx, req.GroupID)
	if err != nil {
		return DepartmentResponse{}, commonErrors.NotFoundError(group.ErrGroupNotFound.Error(), group.ErrGroupNotFound)
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
		return DepartmentResponse{}, commonErrors.NotFoundError(ErrDepartmentNotFound.Error(), ErrDepartmentNotFound)
	}

	if req.Label != dept.Label {
		exisiting, err := s.departmentRepo.FindDepartmentByLabel(ctx, req.Label)
		if err == nil && exisiting != nil {
			return DepartmentResponse{}, commonErrors.ConflictError(ErrDepartmentLabelExists.Error(), ErrDepartmentLabelExists)
		}
	}

	var g *group.Group
	if req.GroupID != dept.GroupID {
		g, err = s.groupRepo.FindGroupByID(ctx, req.GroupID)
		if err != nil {
			return DepartmentResponse{}, commonErrors.NotFoundError(group.ErrGroupNotFound.Error(), group.ErrGroupNotFound)
		}
	} else {
		g, err = s.groupRepo.FindGroupByID(ctx, req.GroupID)
		if err != nil {
			return DepartmentResponse{}, commonErrors.InternalServerError("failed to fetch group", err)
		}
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
		return commonErrors.NotFoundError(ErrDepartmentNotFound.Error(), ErrDepartmentNotFound)
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
