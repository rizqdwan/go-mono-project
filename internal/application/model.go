package application

import (
	"time"
)

type Application struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Year        int64     `db:"year"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdateAt    time.Time `db:"updated_at"`
}

type ApplicationResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Year        int64     `json:"year"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}
