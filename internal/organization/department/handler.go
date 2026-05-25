package department

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

// DepartmentList godoc
// @Summary      Get list of department
// @Description  Get list of all departments with group details
// @Tags         Department
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.WebResponse[DepartmentResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /departments [get]
func (h *Handler) DepartmentList(c *echo.Context) error {
	resp, err := h.svc.ListDepartments(c.Request().Context())
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "department list retrieved successfully", resp)
}

// CreateDepartment godoc
// @Summary 	Create a department
// @Description Create a new department (ProjectAdmin only)
// @Tags 		Department
// @Accept 		json
// @Produce 	json
// @Security    BearerAuth
// @Param        request  body      NewDepartmentRequest  true  "Create Department Request"
// @Success      201  {object}  response.WebResponse[DepartmentResponse]
// @Failure      400  {object}  response.WebResponse[any]
// @Failure      401  {object}  response.WebResponse[any]
// @Failure      404  {object}  response.WebResponse[any]
// @Failure      409  {object}  response.WebResponse[any]
// @Router       /departments [post]
func (h *Handler) CreateDepartment(c *echo.Context) error {
	var req NewDepartmentRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	resp, err := h.svc.CreateDepartment(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusCreated, "department created successfully", resp)
}

// UpdateDepartment godoc
// @Summary      Update a department
// @Description  Update an existing department (ProjectAdmin only)
// @Tags         Department
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                     true  "Department ID"
// @Param        request  body      UpdateDepartmentRequest true  "Update Department Request"
// @Success      200  {object}  response.WebResponse[DepartmentResponse]
// @Failure      400  {object}  response.WebResponse[any]
// @Failure      401  {object}  response.WebResponse[any]
// @Failure      404  {object}  response.WebResponse[any]
// @Failure      409  {object}  response.WebResponse[any]
// @Router       /departments/{id}/details [put]
func (h *Handler) UpdateDepartment(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return response.Error(c, http.StatusBadRequest, "invalid department id", nil)
	}

	var req UpdateDepartmentRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	resp, err := h.svc.UpdateDepartment(c.Request().Context(), id, req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "department updated successfully", resp)
}

// DeleteDepartment godoc
// @Summary      Delete a department
// @Description  Delete a department. Blocked if department has active users. (ProjectAdmin only)
// @Tags         Department
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Department ID"
// @Success      200  {object}  response.WebResponse[any]
// @Failure      400  {object}  response.WebResponse[any]
// @Failure      401  {object}  response.WebResponse[any]
// @Failure      404  {object}  response.WebResponse[any]
// @Failure      409  {object}  response.WebResponse[any]
// @Router       /departments/{id} [delete]
func (h *Handler) DeleteDepartment(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return response.Error(c, http.StatusBadRequest, "invalid department id", nil)
	}

	if err := h.svc.DeleteDepartment(c.Request().Context(), id); err != nil {
		return err
	}

	return response.Success[any](c, http.StatusOK, "department deleted successfully", nil)
}
