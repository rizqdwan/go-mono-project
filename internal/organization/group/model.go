package group

import "time"

type Group struct {
	ID        int64     `db:"id"`
	Label     string    `db:"label"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type GroupResponse struct {
	ID        int64     `json:"id"`
	Label     string    `json:"label"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupInfo struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

type NewGroupRequest struct {
	Label string `json:"label" validate:"required"`
	Name  string `json:"name" validate:"required"`
}

type UpdateGroupRequest struct {
	Label string `json:"label" validate:"required"`
	Name  string `json:"name" validate:"required"`
}
