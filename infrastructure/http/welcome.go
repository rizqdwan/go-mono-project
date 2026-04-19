package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type WelcomeHandler struct{}

func NewWelcomeHandler() *WelcomeHandler {
	return &WelcomeHandler{}
}

// Welcome godoc
// @Summary      Welcome
// @Description  Check if the API is running
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       / [get]
func (h *WelcomeHandler) Welcome(c *echo.Context) error {
	response := map[string]interface{}{
		"message": "Welcome to Go Mono Project API!",
		"status":  "success",
		"code":    http.StatusOK,
	}

	return c.JSON(http.StatusOK, response)
}
