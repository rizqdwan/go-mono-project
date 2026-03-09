package user

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, u *User) error
	UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error
	FindRoleByName(ctx context.Context, name string) (int64, error)
	FindDepartmentByLabel(ctx context.Context, label string) (int64, error)
	FindPositionByName(ctx context.Context, name string) (string, error)
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
		       created_at, updated_at
		FROM users
		WHERE email = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, query, email))
}

func (r *repository) FindByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, name, password_hash, role_id, department_id, position_id,
		       created_at, updated_at
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

func (r *repository) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password_hash = $1,
		    updated_at    = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, hashedPassword, userID)
	return err
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

func (r *repository) FindPositionByName(ctx context.Context, name string) (string, error) {
	query := `
		SELECT id FROM user_positions
		WHERE name = $1
		LIMIT 1
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrPositionNotFound
	}
	return id, err
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
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return u, err
}
