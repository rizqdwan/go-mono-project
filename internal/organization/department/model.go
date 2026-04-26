package department

import "time"

type Department struct {
	ID         int64     `db:"id"`
	Label      string    `db:"label"`
	Name       string    `db:"name"`
	group_id   int64     `db:"group_id"`
	created_at time.Time `db:"created_at"`
	updated_at time.Time `db:"updated_at"`
}

type DepartmentResponse struct {
	ID        int64  `json:"id"`
	Label     string `json:"label"`
	Name      string `json:"name"`
	GroupID   int64  `json:"group_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type DepartmentInfo struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}
