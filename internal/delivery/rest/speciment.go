package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	specimentuc "github.com/oibacidem/lims-hl-seven/internal/usecase/speciment"
)

type SpecimenHandler struct {
	cfg              *config.Schema
	specimentUsecase *specimentuc.SpecimenUseCase
}

func NewSpecimenHandler(cfg *config.Schema, specimentUsecase *specimentuc.SpecimenUseCase) *SpecimenHandler {
	return &SpecimenHandler{cfg: cfg, specimentUsecase: specimentUsecase}
}

func (h *SpecimenHandler) FindSpecimens(c echo.Context) error {
	var req entity.SpecimentGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	speciments, err := h.specimentUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(speciments)))
	return c.JSON(http.StatusOK, speciments)
}

func (h *SpecimenHandler) GetOneSpecimen(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	speciment, err := h.specimentUsecase.FindOneByID(id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, speciment)
}

func (h *SpecimenHandler) DeleteSpecimen(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	if err := h.specimentUsecase.Delete(id); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.Specimen{
		ID: id,
	})
}

func (h *SpecimenHandler) UpdateSpecimen(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req := entity.Specimen{ID: id}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.specimentUsecase.Update(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, req)
}

func (h *SpecimenHandler) CreateSpecimen(c echo.Context) error {
	var req entity.Specimen
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.specimentUsecase.Create(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, req)
}
