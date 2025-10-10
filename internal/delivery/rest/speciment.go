package rest

import (
	"net/http"
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	specimenuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/specimen"
	"github.com/labstack/echo/v4"
)

type SpecimenHandler struct {
	cfg             *config.Schema
	SpecimenUsecase *specimenuc.SpecimenUseCase
}

func NewSpecimenHandler(cfg *config.Schema, SpecimenUsecase *specimenuc.SpecimenUseCase) *SpecimenHandler {
	return &SpecimenHandler{cfg: cfg, SpecimenUsecase: SpecimenUsecase}
}

func (h *SpecimenHandler) FindSpecimens(c echo.Context) error {
	var req entity.SpecimenGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	Specimens, err := h.SpecimenUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, Specimens)
}

func (h *SpecimenHandler) GetOneSpecimen(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	specimen, err := h.SpecimenUsecase.FindOneByID(id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, specimen)
}
