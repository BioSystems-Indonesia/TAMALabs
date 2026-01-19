package rest

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/test_type"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/labstack/echo/v4"
)

type TestTypeHandler struct {
	cfg             *config.Schema
	testTypeUsecase *test_type.Usecase
}

func NewTestTypeHandler(cfg *config.Schema, TestTypeUsecase *test_type.Usecase) *TestTypeHandler {
	return &TestTypeHandler{cfg: cfg, testTypeUsecase: TestTypeUsecase}
}

func (h *TestTypeHandler) ListTestType(c echo.Context) error {
	var req entity.TestTypeGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.testTypeUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, resp)
}

func (h *TestTypeHandler) ListTestTypeFilter(c echo.Context) error {
	var req entity.TestTypeGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.testTypeUsecase.ListAllFilter(
		c.Request().Context(),
	)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *TestTypeHandler) GetOneTestType(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	testType, err := h.testTypeUsecase.FindOneByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}

func (h *TestTypeHandler) GetOneTestTypeByCode(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return handleError(c, entity.ErrBadRequest.WithInternal(nil))
	}

	testType, err := h.testTypeUsecase.FindOneByCode(c.Request().Context(), code)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}

func (h *TestTypeHandler) GetOneTestTypeByAliasCode(c echo.Context) error {
	aliasCode := c.Param("alias_code")
	if aliasCode == "" {
		return handleError(c, entity.ErrBadRequest.WithInternal(nil))
	}

	testType, err := h.testTypeUsecase.FindOneByAliasCode(c.Request().Context(), aliasCode)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}

func (h *TestTypeHandler) CreateTestType(c echo.Context) error {
	var req struct {
		entity.TestType
		DeviceIDs []int `json:"device_ids"`
	}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	// Convert device_ids to Devices array
	if len(req.DeviceIDs) > 0 {
		req.TestType.Devices = make([]entity.Device, len(req.DeviceIDs))
		for i, id := range req.DeviceIDs {
			req.TestType.Devices[i] = entity.Device{ID: id}
		}
	}

	testType, err := h.testTypeUsecase.Create(c.Request().Context(), &req.TestType)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, testType)
}

func (h *TestTypeHandler) UpdateTestType(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	var req struct {
		entity.TestType
		DeviceIDs []int `json:"device_ids"`
	}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	// Log received device_ids
	slog.Info("UpdateTestType received device_ids", "test_type_id", id, "device_ids", req.DeviceIDs, "count", len(req.DeviceIDs))

	req.TestType.ID = id

	// Convert device_ids to Devices array
	if len(req.DeviceIDs) > 0 {
		req.TestType.Devices = make([]entity.Device, len(req.DeviceIDs))
		for i, deviceID := range req.DeviceIDs {
			req.TestType.Devices[i] = entity.Device{ID: deviceID}
		}
		slog.Info("Converted device_ids to Devices", "devices_count", len(req.TestType.Devices))
	} else {
		// If no device_ids provided, clear the devices association
		req.TestType.Devices = []entity.Device{}
		slog.Info("No device_ids provided, clearing devices association")
	}

	testType, err := h.testTypeUsecase.Update(c.Request().Context(), &req.TestType)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}

func (h *TestTypeHandler) DeleteTestType(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req, err := h.testTypeUsecase.FindOneByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	testType, err := h.testTypeUsecase.Delete(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}

func (h *TestTypeHandler) UploadBulkTestType(c echo.Context) error {
	mf, err := c.FormFile("file")
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	f, err := mf.Open()
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}
	defer f.Close()

	err = h.testTypeUsecase.BulkCreate(c.Request().Context(), f)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// TestTypeSimpleResponse is a simplified response for test type list
type TestTypeSimpleResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	Unit           string `json:"unit"`
	ReferenceRange string `json:"reference_range"`
}

func (h *TestTypeHandler) ListAllTestTypes(c echo.Context) error {
	testTypes, err := h.testTypeUsecase.FindAllSimple(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	// Transform to simple response
	response := make([]TestTypeSimpleResponse, len(testTypes))
	for i, tt := range testTypes {
		response[i] = TestTypeSimpleResponse{
			ID:             tt.ID,
			Name:           tt.Name,
			Category:       tt.Category,
			Unit:           tt.Unit,
			ReferenceRange: tt.GetReferenceRange(),
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": response,
	})
}
