package user

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, u *User) error
	UpdatePassword(ctx context.Context, userID int64, hashedPassword string) (time.Time, error)
	FindRoleByName(ctx context.Context, name string) (int64, error)
	FindRoleNameByID(ctx context.Context, roleID int64) (string, error)
	FindDepartmentByLabel(ctx context.Context, label string) (int64, error)
	FindPositionByLabel(ctx context.Context, label string) (string, error)
	FindDepartmentLabelByID(ctx context.Context, id int64) (string, error)
	FindUsersByDepartmentID(ctx context.Context, departmentID int64) ([]UserListResponse, error)
	FindUserDetailsByID(ctx context.Context, userID int64) (*UserDetailsResponse, error)
	DeactivateUser(ctx context.Context, userID int64) error
	HasActiveProjectAssignments(ctx context.Context, userID int64) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
        SELECT id, email, name, password_hash, role_id, department_id, position_id,
               is_active, created_at, updated_at
        FROM users
        WHERE email = $1
    `
	return r.scanUser(r.db.QueryRowContext(ctx, query, email))
}

func (r *repository) FindByID(ctx context.Context, id int64) (*User, error) {
	query := `
        SELECT id, email, name, password_hash, role_id, department_id, position_id,
               is_active, created_at, updated_at
        FROM users
        WHERE id = $1
    `
	return r.scanUser(r.db.QueryRowContext(ctx, query, id))
}

func (r *repository) CreateUser(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users
			(email, name, password_hash, role_id, department_id, position_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		u.Email,
		u.Name,
		u.PasswordHash,
		u.RoleID,
		u.DepartmentID,
		u.PositionID,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *repository) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) (time.Time, error) {
	query := `
		UPDATE users
		SET password_hash = $1,
		    updated_at    = NOW()
		WHERE id = $2
		RETURNING updated_at
	`
	var updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, hashedPassword, userID).Scan(&updatedAt)
	if err != nil {
		return time.Time{}, err
	}
	return updatedAt, nil
}

func (r *repository) FindRoleByName(ctx context.Context, name string) (int64, error) {
	query := `
		SELECT id FROM user_roles
		WHERE name = $1
		LIMIT 1
	`
	var id int64
	err := r.db.QueryRowContext(ctx, query, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrRoleNotFound
	}
	return id, err
}

func (r *repository) FindRoleNameByID(ctx context.Context, roleID int64) (string, error) {
	query := `
        SELECT name FROM user_roles
        WHERE id = $1
        LIMIT 1
    `
	var name string
	err := r.db.QueryRowContext(ctx, query, roleID).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrRoleNotFound
	}
	return name, err
}

func (r *repository) FindDepartmentByLabel(ctx context.Context, label string) (int64, error) {
	query := `
		SELECT id FROM departments
		WHERE label = $1
		LIMIT 1
	`
	var id int64
	err := r.db.QueryRowContext(ctx, query, label).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrDepartmentNotFound
	}
	return id, err
}

func (r *repository) FindPositionByLabel(ctx context.Context, label string) (string, error) {
	query := `
        SELECT id FROM user_positions
        WHERE id = $1
        LIMIT 1
    `
	var id string
	err := r.db.QueryRowContext(ctx, query, label).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrPositionNotFound
	}
	return id, err
}

func (r *repository) FindDepartmentLabelByID(ctx context.Context, id int64) (string, error) {
	query := `
		SELECT label FROM departments
		WHERE id = $1
		LIMIT 1
	`
	var label string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&label)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrDepartmentNotFound
	}
	return label, err
}

func (r *repository) FindUserDetailsByID(ctx context.Context, userID int64) (*UserDetailsResponse, error) {
	query := `
        SELECT u.id, u.email, u.name, ur.name , d.label, d.name, up.name , u.created_at
        FROM users u
        JOIN user_roles ur ON ur.id = u.role_id
        JOIN departments d ON d.id = u.department_id
        JOIN user_positions up ON up.id = u.position_id
        WHERE u.id = $1
        LIMIT 1
      `
	var resp UserDetailsResponse
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&resp.ID,
		&resp.Email,
		&resp.Name,
		&resp.Role,
		&resp.Department.Label,
		&resp.Department.Name,
		&resp.Position,
		&resp.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *repository) FindUsersByDepartmentID(ctx context.Context, departmentID int64) ([]UserListResponse, error) {
	query := `
        SELECT u.id, u.email, u.name, ur.name , d.label, d.name, up.name , u.created_at
        FROM users u
        JOIN user_roles ur ON ur.id = u.role_id
        JOIN departments d ON d.id = u.department_id
        JOIN user_positions up ON up.id = u.position_id
        WHERE u.department_id = $1
    `
	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp []UserListResponse
	for rows.Next() {
		var u UserListResponse
		if err = rows.Scan(
			&u.ID,
			&u.Email,
			&u.Name,
			&u.Role,
			&u.Department.Label,
			&u.Department.Name,
			&u.Position,
			&u.CreatedAt,
		); err != nil {
			return nil, err
		}
		resp = append(resp, u)
	}

	return resp, nil
}

func (r *repository) DeleteUser(ctx context.Context, userID int64) (bool, error) {
	query := `
		DELETE FROM users 
		WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) DeactivateUser(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET is_active  = false,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) HasActiveProjectAssignments(ctx context.Context, userID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM project_members pm
			JOIN projects p ON p.id = pm.project_id
			WHERE pm.user_id = $1
			  AND p.status_id NOT IN ('DONE', 'COMPLETED', 'CANCELLED')
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) scanUser(row *sql.Row) (*User, error) {
	u := &User{}
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.PasswordHash,
		&u.RoleID,
		&u.DepartmentID,
		&u.PositionID,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return u, err
}
