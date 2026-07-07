package group

import (
	"context"
	"errors"

	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
	"github.com/rizqdwan/go-mono-project/pkg/pagination"
)

type Service interface {
	ListGroup(ctx context.Context, p pagination.Params) ([]GroupResponse, int, error)
	CreateGroup(ctx context.Context, req NewGroupRequest) (GroupResponse, error)
	UpdateGroup(ctx context.Context, id int64, req UpdateGroupRequest) (GroupResponse, error)
	DeleteGroup(ctx context.Context, id int64) error
}

type service struct {
	groupRepo Repository
}

func NewService(groupRepo Repository) Service {
	s := &service{
		groupRepo: groupRepo,
	}

	return s
}

func (s *service) ListGroup(ctx context.Context, p pagination.Params) ([]GroupResponse, int, error) {
	group, total, err := s.groupRepo.FindAllGroups(ctx, p)
	if err != nil {
		return nil, 0, commonErrors.InternalServerError("failed to fetch group", err)
	}
	return group, total, nil
}

func (s *service) CreateGroup(ctx context.Context, req NewGroupRequest) (GroupResponse, error) {
	exisiting, err := s.groupRepo.FindGroupByLabel(ctx, req.Label)
	if err == nil && exisiting != nil {
		return GroupResponse{}, commonErrors.ConflictError(ErrGroupLabelExists.Error(), ErrGroupLabelExists)
	}

	if err != nil && !errors.Is(err, ErrGroupNotFound) {
		return GroupResponse{}, commonErrors.InternalServerError("failed to group label", err)
	}

	g := &Group{
		Label: req.Label,
		Name:  req.Label,
	}

	if err := s.groupRepo.CreateGroup(ctx, g); err != nil {
		return GroupResponse{}, commonErrors.InternalServerError("failed to create group", err)
	}

	return GroupResponse{
		ID:        g.ID,
		Label:     g.Label,
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}, nil

}

func (s *service) UpdateGroup(ctx context.Context, id int64, req UpdateGroupRequest) (GroupResponse, error) {
	group, err := s.groupRepo.FindGroupByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrGroupNotFound) {
			return GroupResponse{}, commonErrors.NotFoundError(ErrGroupNotFound.Error(), ErrGroupNotFound)
		}
		return GroupResponse{}, commonErrors.InternalServerError("failed to fetch group", err)
	}

	if req.Label != group.Label {
		existing, err := s.groupRepo.FindGroupByLabel(ctx, req.Label)
		if err == nil && existing != nil {
			return GroupResponse{}, commonErrors.ConflictError(ErrGroupLabelExists.Error(), ErrGroupLabelExists)
		}

		if err != nil && !errors.Is(err, ErrGroupNotFound) {
			return GroupResponse{}, commonErrors.InternalServerError("failed to group label", err)
		}
	}

	group.Label = req.Label
	group.Name = req.Name

	updatedAt, err := s.groupRepo.UpdateGroup(ctx, group)
	if err != nil {
		return GroupResponse{}, commonErrors.InternalServerError("failed to update group", err)
	}

	return GroupResponse{
		ID:        group.ID,
		Label:     group.Label,
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *service) DeleteGroup(ctx context.Context, id int64) error {
	_, err := s.groupRepo.FindGroupByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrGroupNotFound) {
			return commonErrors.NotFoundError(ErrGroupNotFound.Error(), ErrGroupNotFound)
		}
		return commonErrors.InternalServerError("failed to fetch group", err)
	}

	hasDepartment, err := s.groupRepo.HasActiveDepartment(ctx, id)
	if err != nil {
		return commonErrors.InternalServerError("failed to check if group is active", err)
	}
	if hasDepartment {
		return commonErrors.ConflictError(ErrGroupHasActiveDepartment.Error(), ErrGroupHasActiveDepartment)
	}

	if err := s.groupRepo.DeleteGroup(ctx, id); err != nil {
		return commonErrors.InternalServerError("failed to delete group", err)
	}
	return nil
}
