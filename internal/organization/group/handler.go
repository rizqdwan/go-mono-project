package group

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/pkg/pagination"
	"github.com/rizqdwan/go-mono-project/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// GroupList godoc
// @Summary      Get list of group
// @Description  Get list of all group
// @Tags         Group
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.WebResponse[GroupResponse]
// @Failure      400      {object}  response.WebResponse[any]
// @Failure      401      {object}  response.WebResponse[any]
// @Router       /groups [get]
func (h *Handler) GroupList(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	p := pagination.NewParams(page, size)

	groups, total, err := h.svc.ListGroup(c.Request().Context(), p)
	if err != nil {
		return err
	}

	totalPages := (total + p.Size - 1) / p.Size
	meta := response.PaginationMeta(p.Page, p.Size, total, totalPages, p.Page >= totalPages)

	return response.SuccessWithData(c, http.StatusOK, "group list retrieved successfully", groups, meta)
}

// CreateGroup godoc
// @Summary 	Create a group
// @Description Create a new group (ProjectAdmin only)
// @Tags 		Group
// @Accept 		json
// @Produce 	json
// @Security    BearerAuth
// @Param        request  body      NewGroupRequest  true  "Create Group Request"
// @Success      201  {object}  response.WebResponse[GroupResponse]
// @Failure      400  {object}  response.WebResponse[any]
// @Failure      401  {object}  response.WebResponse[any]
// @Failure      404  {object}  response.WebResponse[any]
// @Failure      409  {object}  response.WebResponse[any]
// @Router       /groups [post]
func (h *Handler) CreateGroup(c *echo.Context) error {
	var req NewGroupRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", []string{err.Error()})
	}

	resp, err := h.svc.CreateGroup(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusCreated, "create group successfully", resp)
}

// UpdateGroup godoc
// @Summary      Update a group
// @Description  Update an existing group (ProjectAdmin only)
// @Tags         Group
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                     true  "Group ID"
// @Param        request  body      UpdateGroupRequest true  "Update Group Request"
// @Success      200  {object}  response.WebResponse[GroupResponse]
// @Failure      400  {object}  response.WebResponse[any]
// @Failure      401  {object}  response.WebResponse[any]
// @Failure      404  {object}  response.WebResponse[any]
// @Failure      409  {object}  response.WebResponse[any]
// @Router       /groups/{id}/details [put]
func (h *Handler) UpdateGroup(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return response.Error(c, http.StatusBadRequest, "invalid group id ", nil)
	}

	var req UpdateGroupRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", []string{err.Error()})
	}

	resp, err := h.svc.UpdateGroup(c.Request().Context(), id, req)
	if err != nil {
		return err
	}
	return response.Success(c, http.StatusOK, "update group successfully", resp)
}

// DeleteGroup godoc
// @Summary      Delete a group
// @Description  Delete a group. Blocked if group has active departments. (ProjectAdmin only)
// @Tags         Group
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  response.WebResponse[any]
// @Failure      400  {object}  response.WebResponse[any]
// @Failure      401  {object}  response.WebResponse[any]
// @Failure      404  {object}  response.WebResponse[any]
// @Failure      409  {object}  response.WebResponse[any]
// @Router       /groups/{id} [delete]
func (h *Handler) DeleteGroup(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return response.Error(c, http.StatusBadRequest, "invalid group id ", []string{err.Error()})
	}
	if err := h.svc.DeleteGroup(c.Request().Context(), id); err != nil {
		return err
	}

	return response.Success[any](c, http.StatusOK, "delete group successfully", nil)
}
