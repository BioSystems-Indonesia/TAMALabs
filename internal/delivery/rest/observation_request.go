package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/observation_request"
)

type ObservationRequestHandler struct {
	cfg                       *config.Schema
	ObservationRequestUsecase *observation_requestuc.ObservationRequestUseCase
}

func NewObservationRequestHandler(cfg *config.Schema, ObservationRequestUsecase *observation_requestuc.ObservationRequestUseCase) *ObservationRequestHandler {
	return &ObservationRequestHandler{cfg: cfg, ObservationRequestUsecase: ObservationRequestUsecase}
}

func (h *ObservationRequestHandler) FindObservationRequests(c echo.Context) error {
	var req entity.ObservationRequestGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	ObservationRequests, err := h.ObservationRequestUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(ObservationRequests)))
	return c.JSON(http.StatusOK, ObservationRequests)
}

func (h *ObservationRequestHandler) GetOneObservationRequest(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	observationRequest, err := h.ObservationRequestUsecase.FindOneByID(id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, observationRequest)
}
