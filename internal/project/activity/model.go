package activity

import "time"

type activity struct {
	ID        int64     `db:"id"`
	ProjectID string    `db:"project_id"`
	UserID    int64     `db:"user_id"`
	Progress  int64     `db:"progress"`
	Task      string    `db:"task"`
	StatusID  string    `db:"status_id"`
	CreatedAt time.Time `db:"created_at"`
}
