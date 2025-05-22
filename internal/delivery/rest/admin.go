package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	adminuc "github.com/oibacidem/lims-hl-seven/internal/usecase/admin"
)

type AdminHandler struct {
	cfg          *config.Schema
	adminUsecase *adminuc.AdminUsecase
}

func NewAdminHandler(cfg *config.Schema, adminUsecase *adminuc.AdminUsecase) *AdminHandler {
	return &AdminHandler{cfg: cfg, adminUsecase: adminUsecase}
}

func (h *AdminHandler) FindAdmins(c echo.Context) error {
	var req entity.GetManyRequestAdmin
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.adminUsecase.GetAllAdmin(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, result)
}

func (h *AdminHandler) GetOneAdmin(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	admin, err := h.adminUsecase.GetOneAdmin(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, admin)
}

func (h *AdminHandler) DeleteAdmin(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	if err := h.adminUsecase.DeleteAdmin(c.Request().Context(), id); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.Admin{
		ID: id,
	})
}

func (h *AdminHandler) UpdateAdmin(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req := entity.Admin{ID: id}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.adminUsecase.UpdateAdmin(c.Request().Context(), &req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, req)
}

func (h *AdminHandler) CreateAdmin(c echo.Context) error {
	var req entity.Admin
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.adminUsecase.CreateAdmin(c.Request().Context(), &req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, req)
}
