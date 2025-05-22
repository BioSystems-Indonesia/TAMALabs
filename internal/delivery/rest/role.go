package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	roleuc "github.com/oibacidem/lims-hl-seven/internal/usecase/role"
)

type RoleHandler struct {
	cfg         *config.Schema
	roleUsecase *roleuc.RoleUsecase
}

func NewRoleHandler(cfg *config.Schema, roleUsecase *roleuc.RoleUsecase) *RoleHandler {
	return &RoleHandler{cfg: cfg, roleUsecase: roleUsecase}
}

func (h *RoleHandler) FindRoles(c echo.Context) error {
	var req entity.GetManyRequestRole
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.roleUsecase.GetAllRole(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, result)
}

func (h *RoleHandler) GetOneRole(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	role, err := h.roleUsecase.GetOneRole(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, role)
}
