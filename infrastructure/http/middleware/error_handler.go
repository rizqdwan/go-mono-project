package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	commonErrors "github.com/rizqdwan/go-mono-project/internal/common/errors"
	"github.com/rizqdwan/go-mono-project/pkg/response"
)

func ErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}

	var appErr commonErrors.AppError
	if errors.As(err, &appErr) {
		_ = response.Error(c, appErr.HttpCode, appErr.Message, nil)
		return
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		_ = response.Error(c, echoErr.Code, echoErr.Message, nil)
		return
	}

	_ = response.Error(c, http.StatusInternalServerError, "internal server error", nil)
}
