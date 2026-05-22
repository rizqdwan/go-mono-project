package auth

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/infrastructure/http/middleware"
	"github.com/rizqdwan/go-mono-project/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account (ProjectAdmin only)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      RegisterRequest  true  "Register Request"
// @Success      201      {object}  response.WebResponse[RegisterResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /auth/signup [post]
func (h *Handler) Register(c *echo.Context) error {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	resp, err := h.svc.Register(c.Request().Context(), adminID, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, "user registered successfully", resp)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return access + refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login Request"
// @Success      200      {object}  response.WebResponse[TokenResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /auth/signin [post]
func (h *Handler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	browserInfo := c.Request().Header.Get("User-Agent")

	resp, err := h.svc.Login(c.Request().Context(), req, browserInfo)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "login successful", resp)
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate the provided refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      RefreshTokenRequest  true  "Refresh Token"
// @Success      200      {object}  response.WebResponse[any]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	if err := h.svc.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return err
	}

	return response.Success[any](c, http.StatusOK, "logged out successfully", nil)
}

// RenewToken godoc
// @Summary      Renew access token
// @Description  Generate a new access token using a valid refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      RefreshTokenRequest  true  "Refresh Token"
// @Success      200      {object}  response.WebResponse[TokenResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /auth/renew [post]
func (h *Handler) RenewToken(c *echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	browserInfo := c.Request().Header.Get("User-Agent")

	resp, err := h.svc.RenewToken(c.Request().Context(), req, browserInfo)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "token renewed successfully", resp)
}
