package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// UserDetails godoc
// @Summary      Get user details
// @Description  Get the authentication user details
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      201      {object}  response.WebResponse[UserDetailsResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /auth/signup [post]
func (h *Handler) UserDetails(c *echo.Context) error {
	UserID, ok := c.Get("UserID").(int64)
	if !ok || UserID == 0 {
		return response.Error(c, http.StatusUnauthorized, "missing user identity", nil)
	}

	resp, err := h.svc.UserDetails(c.Request().Context(), UserID)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, "User details successfully", resp)
}

// ChangePassword godoc
// @Summary      Change own password
// @Description  Authenticated user changes their own password. Requires current password.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      ChangePasswordRequest  true  "Change Password Request"
// @Success      200      {object}  response.WebResponse[ChangePasswordResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /user/change-password [put]
func (h *Handler) ChangePassword(c *echo.Context) error {
	userID, ok := c.Get("userID").(int64)
	if !ok || userID == 0 {
		return response.Error(c, http.StatusUnauthorized, "missing user identity", nil)
	}

	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	resp, err := h.svc.ChangePassword(c.Request().Context(), userID, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "password changed successfully", resp)
}

// ResetPassword godoc
// @Summary      Reset a user's password (Admin)
// @Description  ProjectAdmin resets any user's password. Does not require the current password. Forces logout on all active sessions.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                    true  "Target User ID"
// @Param        request  body      ResetPasswordRequest   true  "Reset Password Request"
// @Success      200      {object}  response.WebResponse[ResetPasswordResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Failure      403      {object}  response.WebResponse[any]
// @Failure      404      {object}  response.WebResponse[any]
// @Router       /user/{id}/reset-password [put]
func (h *Handler) ResetPassword(c *echo.Context) error {
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetUserID == 0 {
		return response.Error(c, http.StatusBadRequest, "invalid user id", nil)
	}

	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	resp, err := h.svc.ResetPassword(c.Request().Context(), targetUserID, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "password reset successfully", resp)
}
