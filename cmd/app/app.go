package main

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/config"
	_ "github.com/rizqdwan/go-mono-project/docs"
	"github.com/rizqdwan/go-mono-project/infrastructure/db"
	infrahttp "github.com/rizqdwan/go-mono-project/infrastructure/http"
	"github.com/rizqdwan/go-mono-project/infrastructure/http/middleware"
	"github.com/rizqdwan/go-mono-project/infrastructure/security"
	"github.com/rizqdwan/go-mono-project/internal/auth"
	"github.com/rizqdwan/go-mono-project/internal/organization/department"
	"github.com/rizqdwan/go-mono-project/internal/organization/group"
	"github.com/rizqdwan/go-mono-project/internal/user"
	"github.com/rizqdwan/go-mono-project/pkg/token"
	"github.com/rizqdwan/go-mono-project/pkg/validator"
)

type application struct {
	echo     *echo.Echo
	database *db.Database
}

func newApplication(ctx context.Context, cfg *config.Config) (*application, error) {
	database, err := db.NewDatabase(cfg.DB)
	if err != nil {
		return nil, err
	}

	tokenSvc := token.NewJWTService(cfg.JWT)
	passwordSvc := security.NewPasswordHash()

	userRepo := user.NewRepository(database.DB)
	authRepo := auth.NewRepository(database.DB)
	departmentRepo := department.NewRepository(database.DB)
	groupRepo := group.NewRepository(database.DB)

	authSvc := auth.NewService(ctx, authRepo, userRepo, tokenSvc, passwordSvc)
	userSvc := user.NewService(userRepo, tokenSvc, passwordSvc, authRepo, cfg.Role)
	deptSvc := department.NewService(departmentRepo, groupRepo)
	groupSvc := group.NewService(groupRepo)

	e := echo.New()
	e.Validator = validator.New()
	e.HTTPErrorHandler = middleware.ErrorHandler
	e.Use(middleware.RequestIDMiddleware())
	e.Use(middleware.LoggingMiddleware())

	handlers := &infrahttp.Handlers{
		Auth:       auth.NewHandler(authSvc),
		User:       user.NewHandler(userSvc),
		Department: department.NewHandler(deptSvc),
		Group:      group.NewHandler(groupSvc),
	}

	infrahttp.SetupRouter(e, tokenSvc, cfg.Role, handlers)

	return &application{
		echo:     e,
		database: database,
	}, nil
}

func (a *application) close() {
	err := a.database.Close()
	if err != nil {
		return
	}
}
