package department

import (
	"time"

	"github.com/rizqdwan/go-mono-project/internal/organization/group"
)

type Department struct {
	ID         int64     `db:"id"`
	Label      string    `db:"label"`
	Name       string    `db:"name"`
	Group_id   int64     `db:"group_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
}

// @name DepartmentResponse
type DepartmentResponse struct {
	ID        int64           `json:"id"`
	Label     string          `json:"label"`
	Name      string          `json:"name"`
	Group     group.GroupInfo `json:"group"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type DepartmentInfo struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

type NewDepartmentRequest struct {
	Label   string `json:"label" validate:"required"`
	Name    string `json:"name" validate:"required"`
	GroupID int64  `json:"group_id" validate:"required"`
}

type UpdateDepartmentRequest struct {
	Label   string `json:"label" validate:"required"`
	Name    string `json:"name" validate:"required"`
	GroupID int64  `json:"group_id" validate:"required"`
}
