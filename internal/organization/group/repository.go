package group

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rizqdwan/go-mono-project/pkg/pagination"
)

type Repository interface {
	FindAllGroups(ctx context.Context, p pagination.Params) ([]GroupResponse, int, error)
	FindGroupByID(ctx context.Context, id int64) (*Group, error)
	FindGroupByLabel(ctx context.Context, label string) (*Group, error)
	CreateGroup(ctx context.Context, g *Group) error
	UpdateGroup(ctx context.Context, g *Group) (time.Time, error)
	DeleteGroup(ctx context.Context, id int64) error
	HasActiveDepartment(ctx context.Context, groupID int64) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAllGroups(ctx context.Context, p pagination.Params) ([]GroupResponse, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM groups`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, label, name, created_at, updated_at FROM groups
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
    `
	rows, err := r.db.QueryContext(ctx, query, p.Size, p.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var groups []GroupResponse
	for rows.Next() {
		var group GroupResponse
		if err := rows.Scan(
			&group.ID,
			&group.Label,
			&group.Name,
			&group.CreatedAt,
			&group.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		groups = append(groups, group)
	}
	return groups, total, rows.Err()
}

func (r *repository) FindGroupByID(ctx context.Context, id int64) (*Group, error) {
	query := `
			SELECT id, label, name, created_at, updated_at FROM groups
			WHERE id = $1
			LIMIT 1
    `
	var group Group
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID,
		&group.Label,
		&group.Name,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrGroupNotFound
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *repository) FindGroupByLabel(ctx context.Context, label string) (*Group, error) {
	query := `
		SELECT id, label, name, created_at, updated_at FROM groups
		WHERE label = $1
		LIMIT 1
    `

	var group Group
	err := r.db.QueryRowContext(ctx, query, label).Scan(
		&group.ID,
		&group.Label,
		&group.Name,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrGroupNotFound
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *repository) CreateGroup(ctx context.Context, g *Group) error {
	query := `
		INSERT INTO groups (label, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
    `
	return r.db.QueryRowContext(ctx, query,
		g.Label,
		g.Name,
	).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
}

func (r *repository) UpdateGroup(ctx context.Context, g *Group) (time.Time, error) {
	query := `
		UPDATE groups
        SET label = $2, 
            name = $3,
            updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
    `

	var updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query,
		g.ID,
		g.Label,
		g.Name,
	).Scan(&updatedAt)
	if err != nil {
		return time.Time{}, err
	}
	return updatedAt, nil
}

func (r *repository) DeleteGroup(ctx context.Context, id int64) error {
	query := `DELETE FROM groups WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) HasActiveDepartment(ctx context.Context, groupID int64) (bool, error) {
	query := `
		SELECT EXISTS (
		SELECT 1 FROM departments 
		WHERE group_id = $1 
                     )
    `
	var exists bool
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(&exists)
	return exists, err
}
