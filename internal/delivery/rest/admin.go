package rest

import (
	"net/http"
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	adminuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/admin"
	"github.com/labstack/echo/v4"
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

// ListAllDoctors returns all doctors without pagination
func (h *AdminHandler) ListAllDoctors(c echo.Context) error {
	doctors, err := h.adminUsecase.GetAllDoctors(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	// Map to simplified response
	type DoctorResponse struct {
		ID       int64  `json:"id"`
		Fullname string `json:"fullname"`
		Username string `json:"username"`
		IsActive bool   `json:"is_active"`
	}

	response := make([]DoctorResponse, len(doctors))
	for i, doctor := range doctors {
		response[i] = DoctorResponse{
			ID:       doctor.ID,
			Fullname: doctor.Fullname,
			Username: doctor.Username,
			IsActive: doctor.IsActive,
		}
	}

	return c.JSON(http.StatusOK, response)
}

// ListAllAnalyzers returns all analyzers without pagination
func (h *AdminHandler) ListAllAnalyzers(c echo.Context) error {
	analyzers, err := h.adminUsecase.GetAllAnalyzers(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	// Map to simplified response
	type AnalyzerResponse struct {
		ID       int64  `json:"id"`
		Fullname string `json:"fullname"`
		Username string `json:"username"`
		IsActive bool   `json:"is_active"`
	}

	response := make([]AnalyzerResponse, len(analyzers))
	for i, analyzer := range analyzers {
		response[i] = AnalyzerResponse{
			ID:       analyzer.ID,
			Fullname: analyzer.Fullname,
			Username: analyzer.Username,
			IsActive: analyzer.IsActive,
		}
	}

	return c.JSON(http.StatusOK, response)
}
