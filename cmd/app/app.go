package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/config"
	_ "github.com/rizqdwan/go-mono-project/docs"
	"github.com/rizqdwan/go-mono-project/infrastructure/db"
	infrahttp "github.com/rizqdwan/go-mono-project/infrastructure/http"
	"github.com/rizqdwan/go-mono-project/infrastructure/security"
	"github.com/rizqdwan/go-mono-project/internal/auth"
	"github.com/rizqdwan/go-mono-project/internal/user"
	"github.com/rizqdwan/go-mono-project/pkg/response"
	"github.com/rizqdwan/go-mono-project/pkg/token"
	"github.com/rizqdwan/go-mono-project/pkg/validator"
)

type application struct {
	echo     *echo.Echo
	database *db.Database
}

func newApplication(cfg *config.Config) (*application, error) {
	database, err := db.NewDatabase(cfg.DB)
	if err != nil {
		return nil, err
	}

	tokenSvc := token.NewJWTService(cfg.JWT)
	passwordSvc := security.NewPasswordHash()

	userRepo := user.NewRepository(database.DB)
	authRepo := auth.NewRepository(database.DB)

	authSvc := auth.NewService(authRepo, userRepo, tokenSvc, passwordSvc)
	authHandler := auth.NewHandler(authSvc)

	e := echo.New()
	e.Validator = validator.New()

	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		var he *echo.HTTPError
		if errors.As(err, &he) {
			err := response.Error(c, he.Code, fmt.Sprintf("%v", he.Message), nil)
			if err != nil {
				return
			}
			return
		}
		err = response.Error(c, http.StatusInternalServerError, "internal server error", nil)
		if err != nil {
			return
		}
	}

	infrahttp.SetupRouter(e, tokenSvc, authHandler)

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
