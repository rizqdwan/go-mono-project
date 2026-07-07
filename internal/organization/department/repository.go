package department

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rizqdwan/go-mono-project/pkg/pagination"
)

type Repository interface {
	FindAllDepartments(ctx context.Context, p pagination.Params) ([]DepartmentResponse, int, error)
	FindDepartmentByID(ctx context.Context, id int64) (*Department, error)
	FindDepartmentByLabel(ctx context.Context, label string) (*Department, error)
	FindDepartmentsByGroupID(ctx context.Context, groupID int64) ([]DepartmentResponse, error)
	CreateDepartment(ctx context.Context, d *Department) error
	UpdateDepartment(ctx context.Context, d *Department) (time.Time, error)
	DeleteDepartment(ctx context.Context, id int64) error
	HasActiveUsers(ctx context.Context, departmentID int64) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAllDepartments(ctx context.Context, p pagination.Params) ([]DepartmentResponse, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM departments`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT d.id, d.label, d.name, g.label, g.name, d.created_at, d.updated_at 
			FROM departments d
			JOIN groups g ON g.id = d.group_id
		ORDER BY d.created_at DESC
		LIMIT $1 OFFSET $2
    `
	rows, err := r.db.QueryContext(ctx, query, p.Size, p.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var resp []DepartmentResponse
	for rows.Next() {
		var d DepartmentResponse
		if err = rows.Scan(
			&d.ID,
			&d.Label,
			&d.Name,
			&d.Group.Label,
			&d.Group.Name,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		resp = append(resp, d)
	}
	return resp, total, rows.Err()

}

func (r *repository) FindDepartmentByID(ctx context.Context, id int64) (*Department, error) {
	query := `
			SELECT id, label, name, group_id, created_at, updated_at
        	FROM departments
			WHERE id = $1
			LIMIT 1
    `

	var d Department
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID,
		&d.Label,
		&d.Name,
		&d.GroupID,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDepartmentNotFound
	}
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *repository) FindDepartmentByLabel(ctx context.Context, label string) (*Department, error) {
	query := `
		SELECT id, label, name, group_id, created_at, updated_at
        FROM departments
        WHERE label = $1
        LIMIT 1
	`
	var d Department
	err := r.db.QueryRowContext(ctx, query, label).Scan(
		&d.ID,
		&d.Label,
		&d.Name,
		&d.GroupID,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDepartmentNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repository) FindDepartmentsByGroupID(ctx context.Context, groupID int64) ([]DepartmentResponse, error) {
	query := `
			SELECT d.id, d.label, d.name, g.label, g.name, d.created_at, d.updated_at 
			FROM departments d
			JOIN groups g ON g.id = d.group_id
			WHERE d.group_id = $1
    `
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp []DepartmentResponse
	for rows.Next() {
		var d DepartmentResponse
		if err = rows.Scan(
			&d.ID,
			&d.Label,
			&d.Name,
			&d.Group.Label,
			&d.Group.Name,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		resp = append(resp, d)
	}
	return resp, rows.Err()
}

func (r *repository) CreateDepartment(ctx context.Context, d *Department) error {
	query := `
		INSERT INTO departments
		(label, name, group_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
    `

	return r.db.QueryRowContext(ctx, query,
		d.Label,
		d.Name,
		d.GroupID,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *repository) UpdateDepartment(ctx context.Context, d *Department) (time.Time, error) {
	query := `
		UPDATE departments
        SET label = $2, 
            name = $3, 
            group_id = $4, 
            updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
    `
	var updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query,
		d.ID,
		d.Label,
		d.Name,
		d.GroupID,
	).Scan(&updatedAt)
	if err != nil {
		return time.Time{}, err
	}
	return updatedAt, nil
}

func (r *repository) DeleteDepartment(ctx context.Context, id int64) error {
	query := `DELETE FROM departments WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) HasActiveUsers(ctx context.Context, departmentID int64) (bool, error) {
	query := `
		SELECT EXISTS (
    SELECT 1 FROM users
    WHERE department_id = $1
    AND is_active = true
		)
    `
	var exists bool
	err := r.db.QueryRowContext(ctx, query, departmentID).Scan(&exists)
	return exists, err
}
