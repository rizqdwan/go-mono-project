package application

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	ListApplications(ctx context.Context) ([]ApplicationResponse, error)
	FindApplicationByID(ctx context.Context, id int64) (*Application, error)
	FindApplicationByYear(ctx context.Context, year int64) ([]ApplicationResponse, error)
	CreateApplication(ctx context.Context, app *Application) error
	UpdateApplication(ctx context.Context, app *Application) (time.Time, error)
	DeleteApplication(ctx context.Context, id int64) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
