package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
)

type PatientHandler struct {
	cfg            *config.Schema
	patientUsecase *patientuc.PatientUseCase
}

func NewPatientHandler(cfg *config.Schema, patientUsecase *patientuc.PatientUseCase) *PatientHandler {
	return &PatientHandler{cfg: cfg, patientUsecase: patientUsecase}
}

func (h *PatientHandler) FindPatients(c echo.Context) error {
	var req entity.GetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	patients, err := h.patientUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(patients)))
	return c.JSON(http.StatusOK, patients)
}

func (h *PatientHandler) GetOnePatient(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	patient, err := h.patientUsecase.FindOneByID(id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) DeletePatient(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	if err := h.patientUsecase.Delete(id); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.Patient{
		ID: id,
	})
}

func (h *PatientHandler) UpdatePatient(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req := entity.Patient{ID: id}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.patientUsecase.Update(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, req)
}

func (h *PatientHandler) CreatePatient(c echo.Context) error {
	var req entity.Patient
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.patientUsecase.Create(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, req)
}
