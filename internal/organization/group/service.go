package group

import (
	"context"

	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
)

type Service interface {
	ListGroup(ctx context.Context) ([]GroupResponse, error)
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

func (s *service) ListGroup(ctx context.Context) ([]GroupResponse, error) {
	group, err := s.groupRepo.FindAllGroups(ctx)
	if err != nil {
		return nil, commonErrors.InternalServerError("failed to fetch group", err)
	}
	return group, nil
}

func (s *service) CreateGroup(ctx context.Context, req NewGroupRequest) (GroupResponse, error) {
	exisiting, err := s.groupRepo.FindGroupByLabel(ctx, req.Label)
	if err == nil && exisiting != nil {
		return GroupResponse{}, commonErrors.ConflictError(ErrGroupLabelExists.Error(), ErrGroupLabelExists)
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
		return GroupResponse{}, commonErrors.NotFoundError(ErrGroupNotFound.Error(), ErrGroupNotFound)
	}

	if req.Label != group.Label {
		exisiting, err := s.groupRepo.FindGroupByLabel(ctx, req.Label)
		if err == nil && exisiting != nil {
			return GroupResponse{}, commonErrors.ConflictError(ErrGroupLabelExists.Error(), ErrGroupLabelExists)
		}
	}

	group.Label = req.Label
	group.Name = req.Label

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
		return commonErrors.NotFoundError(ErrGroupNotFound.Error(), ErrGroupNotFound)
	}

	hasDepartment, err := s.groupRepo.HasActiveDepartment(ctx, id)
	if err != nil {
		return commonErrors.InternalServerError("failed to check if group is active", ErrGroupNotFound)
	}
	if hasDepartment {
		return commonErrors.ConflictError(ErrGroupHasActiveDepartment.Error(), ErrGroupHasActiveDepartment)
	}

	if err := s.groupRepo.DeleteGroup(ctx, id); err != nil {
		return commonErrors.InternalServerError("failed to delete group", err)
	}
	return nil
}
