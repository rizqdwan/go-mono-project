package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	FindSessionByRefreshToken(ctx context.Context, refreshToken string) (*Authentication, error)
	FindActiveSessionsByUserID(ctx context.Context, userID int64) ([]Authentication, error)
	FindSessionByUserAndBrowser(ctx context.Context, userID int64, browserInfo string) (*Authentication, error)
	CreateSession(ctx context.Context, session *Authentication) error
	UpdateSession(ctx context.Context, session *Authentication) error
	DeactivateSession(ctx context.Context, id int64) error
	DeactivateAllUserSessions(ctx context.Context, userID int64) error
	FindExpiredActiveSessions(ctx context.Context, cutoff time.Time, limit int) ([]Authentication, error)
	DeleteInactiveSessionsOlderThan(ctx context.Context, before time.Time) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindSessionByRefreshToken(ctx context.Context, refreshToken string) (*Authentication, error) {
	query := `
		SELECT id, user_id, refresh_token, is_active, browser_info, created_at, last_activity
		FROM authentication_sessions
		WHERE refresh_token = $1
	`
	return r.scanSession(r.db.QueryRowContext(ctx, query, refreshToken))
}

func (r *repository) FindActiveSessionsByUserID(ctx context.Context, userID int64) ([]Authentication, error) {
	query := `
		SELECT id, user_id, refresh_token, is_active, browser_info, created_at, last_activity
		FROM authentication_sessions
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Authentication
	for rows.Next() {
		var s Authentication
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.RefreshToken,
			&s.IsActive, &s.BrowserInfo, &s.CreatedAt, &s.LastActivity,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func (r *repository) FindSessionByUserAndBrowser(ctx context.Context, userID int64, browserInfo string) (*Authentication, error) {
	query := `
		SELECT id, user_id, refresh_token, is_active, browser_info, created_at, last_activity
		FROM authentication_sessions
		WHERE user_id = $1
		  AND browser_info = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	return r.scanSession(r.db.QueryRowContext(ctx, query, userID, browserInfo))
}

func (r *repository) CreateSession(ctx context.Context, session *Authentication) error {
	query := `
		INSERT INTO authentication_sessions
			(user_id, refresh_token, is_active, browser_info, created_at, last_activity)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, last_activity
	`
	return r.db.QueryRowContext(ctx, query,
		session.UserID,
		session.RefreshToken,
		session.IsActive,
		session.BrowserInfo,
	).Scan(&session.ID, &session.CreatedAt, &session.LastActivity)
}

func (r *repository) UpdateSession(ctx context.Context, session *Authentication) error {
	query := `
		UPDATE authentication_sessions
		SET refresh_token = $1,
		    is_active     = $2,
		    last_activity = NOW()
		WHERE id = $3
		RETURNING last_activity
	`
	return r.db.QueryRowContext(ctx, query,
		session.RefreshToken,
		session.IsActive,
		session.ID,
	).Scan(&session.LastActivity)
}

func (r *repository) DeactivateSession(ctx context.Context, id int64) error {
	query := `
		UPDATE authentication_sessions
		SET is_active    = false,
		    last_activity = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) DeactivateAllUserSessions(ctx context.Context, userID int64) error {
	query := `
		UPDATE authentication_sessions
		SET is_active    = false,
		    last_activity = NOW()
		WHERE user_id = $1 AND is_active = true
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) FindExpiredActiveSessions(ctx context.Context, cutoff time.Time, limit int) ([]Authentication, error) {
	query := `
		SELECT id, user_id, refresh_token, is_active, browser_info, created_at, last_activity
		FROM authentication_sessions
		WHERE is_active = true
		  AND last_activity < $1
		ORDER BY last_activity ASC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, cutoff, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Authentication
	for rows.Next() {
		var s Authentication
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.RefreshToken,
			&s.IsActive, &s.BrowserInfo, &s.CreatedAt, &s.LastActivity,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func (r *repository) DeleteInactiveSessionsOlderThan(ctx context.Context, before time.Time) error {
	query := `
		DELETE FROM authentication_sessions
		WHERE is_active = false
		  AND last_activity < $1
	`
	_, err := r.db.ExecContext(ctx, query, before)
	return err
}

func (r *repository) scanSession(row *sql.Row) (*Authentication, error) {
	s := &Authentication{}
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshToken,
		&s.IsActive,
		&s.BrowserInfo,
		&s.CreatedAt,
		&s.LastActivity,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSessionNotFound
	}
	return s, err
}
