package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	unitUC "github.com/oibacidem/lims-hl-seven/internal/usecase/unit"
)

type UnitHandler struct {
	cfg    *config.Schema
	unitUC *unitUC.UnitUseCase
}

func NewUnitHandler(cfg *config.Schema, unitUsecase *unitUC.UnitUseCase) *UnitHandler {
	return &UnitHandler{
		cfg:    cfg,
		unitUC: unitUsecase,
	}
}

func (h *UnitHandler) ListUnit(c echo.Context) error {
	var req entity.UnitGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.unitUC.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, result)
}
