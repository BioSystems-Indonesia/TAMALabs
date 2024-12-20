package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/specimen"
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

	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(Specimens)))
	return c.JSON(http.StatusOK, Specimens)
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

func (h *SpecimenHandler) DeleteSpecimen(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	if err := h.SpecimenUsecase.Delete(int(id)); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.Specimen{
		ID: int(id),
	})
}

func (h *SpecimenHandler) UpdateSpecimen(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req := entity.Specimen{ID: int(id)}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.SpecimenUsecase.Update(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, req)
}

func (h *SpecimenHandler) CreateSpecimen(c echo.Context) error {
	var req entity.Specimen
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.SpecimenUsecase.Create(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, req)
}
