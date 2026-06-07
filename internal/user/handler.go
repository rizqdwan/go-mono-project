package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/infrastructure/http/middleware"
	"github.com/rizqdwan/go-mono-project/pkg/pagination"
	"github.com/rizqdwan/go-mono-project/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// UserList godoc
// @Summary      Get list user by department
// @Description  Get list of user with details by department
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.WebResponse[UserListResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /users/list [get]
func (h *Handler) UserList(c *echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	p := pagination.NewParams(page, size)

	users, total, err := h.svc.ListUser(c.Request().Context(), userID, p)
	if err != nil {
		return err
	}

	totalPages := (total + p.Size - 1) / p.Size
	meta := response.PaginationMeta(p.Page, p.Size, total, totalPages, p.Page >= totalPages)

	return response.SuccessWithData(c, http.StatusOK, "user list retrieved successfully", users, meta)
}

// UserDetails godoc
// @Summary      Get current user details
// @Description  Get the authentication user details
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.WebResponse[UserDetailsResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /users/current [get]
func (h *Handler) UserDetails(c *echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	resp, err := h.svc.UserDetails(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, "user details retrieved successfully", resp)
}

// UpdateUserDetails godoc
// @Summary 	Update user details
// @Description Update user details data by project admin
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                    true  "Target User ID"
// @Param        request  body      UpdateUserDetailsRequest   true  "Update User Details Request"
// @Success      200      {object}  response.WebResponse[UserDetailsResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Failure      403      {object}  response.WebResponse[any]
// @Failure      404      {object}  response.WebResponse[any]
// @Router       /users/{id}/details [put]
func (h *Handler) UpdateUserDetails(c *echo.Context) error {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetUserID == 0 {
		return response.Error(c, http.StatusBadRequest, "invalid user id", nil)
	}

	var req UpdateUserDetailsRequest

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	resp, err := h.svc.UpdateUserDetails(c.Request().Context(), adminID, targetUserID, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "user details updated successfully", resp)
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
// @Router       /users/change-password [put]
func (h *Handler) ChangePassword(c *echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
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
// @Router       /users/{id}/reset-password [put]
func (h *Handler) ResetPassword(c *echo.Context) error {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

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

	resp, err := h.svc.ResetPassword(c.Request().Context(), adminID, targetUserID, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "password reset successfully", resp)
}

// DeleteUser godoc
// @Summary      Delete a user (Admin)
// @Description  ProjectAdmin soft-deletes a user from their department. Blocks if user has active project assignments.
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Target User ID"
// @Success      200      {object}  response.WebResponse[any]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Failure      403      {object}  response.WebResponse[any]
// @Failure      404      {object}  response.WebResponse[any]
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUser(c *echo.Context) error {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetUserID == 0 {
		return response.Error(c, http.StatusBadRequest, "invalid user id", nil)
	}

	if err := h.svc.DeleteUser(c.Request().Context(), adminID, targetUserID); err != nil {
		return err
	}

	return response.Success[any](c, http.StatusOK, "user deleted successfully", nil)
}
