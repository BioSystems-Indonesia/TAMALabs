package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	observation_requestuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/observation_request"
	"github.com/labstack/echo/v4"
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

	resp, err := h.ObservationRequestUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, resp)
}

func (h *ObservationRequestHandler) GetOneObservationRequest(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	observationRequest, err := h.ObservationRequestUsecase.FindOneByID(context.TODO(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, observationRequest)
}
