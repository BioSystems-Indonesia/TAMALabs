package rest

import (
	"net/http"
	"strconv"

	externalEntity "github.com/BioSystems-Indonesia/TAMALabs/internal/entity/external"
	appMiddleware "github.com/BioSystems-Indonesia/TAMALabs/internal/middleware"
	technomedicuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/technomedic"
	"github.com/labstack/echo/v4"
)

// TechnoMedicHandler handles HTTP requests for TechnoMedic integration
type TechnoMedicHandler struct {
	usecase               *technomedicuc.Usecase
	integrationMiddleware *appMiddleware.IntegrationCheckMiddleware
}

// NewTechnoMedicHandler creates a new TechnoMedicHandler instance
func NewTechnoMedicHandler(
	usecase *technomedicuc.Usecase,
	integrationMiddleware *appMiddleware.IntegrationCheckMiddleware,
) *TechnoMedicHandler {
	return &TechnoMedicHandler{
		usecase:               usecase,
		integrationMiddleware: integrationMiddleware,
	}
}

// RegisterRoutes registers all TechnoMedic routes with integration check middleware
func (h *TechnoMedicHandler) RegisterRoutes(router *echo.Group) {
	technomedic := router.Group("/technomedic", h.integrationMiddleware.CheckTechnoMedicEnabled())

	// GET endpoints
	technomedic.GET("/test-types", h.GetTestTypes)
	technomedic.GET("/sub-categories", h.GetSubCategories)
	technomedic.GET("/sub-categories/:id/test-types", h.GetTestTypesBySubCategory)
	technomedic.GET("/doctors", h.GetDoctors)
	technomedic.GET("/analysts", h.GetAnalysts)
	technomedic.GET("/order/:no_order", h.GetOrder)

	// POST endpoints
	technomedic.POST("/order", h.CreateOrder)
}

// GetTestTypes returns all available test types
// @Summary Get test types
// @Description Get all available test types for TechnoMedic
// @Tags TechnoMedic
// @Accept json
// @Produce json
// @Success 200 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 500 {object} externalEntity.TechnoMedicOrderResponse
// @Router /technomedic/test-types [get]
func (h *TechnoMedicHandler) GetTestTypes(c echo.Context) error {
	testTypes, err := h.usecase.GetTestTypes(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, externalEntity.TechnoMedicOrderResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Test types retrieved successfully",
		Data:    testTypes,
	})
}

// GetSubCategories returns all unique sub-categories
// @Summary Get sub-categories
// @Description Get all unique sub-categories for TechnoMedic
// @Tags TechnoMedic
// @Accept json
// @Produce json
// @Success 200 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 500 {object} externalEntity.TechnoMedicOrderResponse
// @Router /technomedic/sub-categories [get]
func (h *TechnoMedicHandler) GetSubCategories(c echo.Context) error {
	subCategories, err := h.usecase.GetSubCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, externalEntity.TechnoMedicOrderResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Sub-categories retrieved successfully",
		Data:    subCategories,
	})
}

// GetTestTypesBySubCategory returns all test types for a specific sub-category
// @Summary Get test types by sub-category
// @Description Get all test types for a specific sub-category
// @Tags TechnoMedic
// @Accept json
// @Produce json
// @Param id path int true "Sub-category ID"
// @Success 200 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 400 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 500 {object} externalEntity.TechnoMedicOrderResponse
// @Router /technomedic/sub-categories/{id}/test-types [get]
func (h *TechnoMedicHandler) GetTestTypesBySubCategory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid sub-category ID",
		})
	}

	testTypes, err := h.usecase.GetTestTypesBySubCategory(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, externalEntity.TechnoMedicOrderResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Test types retrieved successfully",
		Data:    testTypes,
	})
}

// GetDoctors returns all doctors
// @Summary Get doctors
// @Description Get all doctors for TechnoMedic
// @Tags TechnoMedic
// @Accept json
// @Produce json
// @Success 200 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 500 {object} externalEntity.TechnoMedicOrderResponse
// @Router /technomedic/doctors [get]
func (h *TechnoMedicHandler) GetDoctors(c echo.Context) error {
	doctors, err := h.usecase.GetDoctors(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, externalEntity.TechnoMedicOrderResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Doctors retrieved successfully",
		Data:    doctors,
	})
}

// GetAnalysts returns all analysts
// @Summary Get analysts
// @Description Get all analysts/analyzers for TechnoMedic
// @Tags TechnoMedic
// @Accept json
// @Produce json
// @Success 200 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 500 {object} externalEntity.TechnoMedicOrderResponse
// @Router /technomedic/analysts [get]
func (h *TechnoMedicHandler) GetAnalysts(c echo.Context) error {
	analysts, err := h.usecase.GetAnalysts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, externalEntity.TechnoMedicOrderResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Analysts retrieved successfully",
		Data:    analysts,
	})
}

// CreateOrder creates a new order from TechnoMedic
// @Summary Create order
// @Description Create a new order from TechnoMedic
// @Tags TechnoMedic
// @Accept json
// @Produce json
// @Param order body externalEntity.TechnoMedicOrderRequest true "Order request"
// @Success 201 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 400 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 500 {object} externalEntity.TechnoMedicOrderResponse
// @Router /technomedic/order [post]
func (h *TechnoMedicHandler) CreateOrder(c echo.Context) error {
	var req externalEntity.TechnoMedicOrderRequest

	if err := bindAndValidate(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	err := h.usecase.CreateOrder(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, externalEntity.TechnoMedicOrderResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Order created successfully",
		Data: map[string]string{
			"no_order": req.NoOrder,
		},
	})
}

// GetOrder retrieves order details including results
// @Summary Get order
// @Description Get order details including results
// @Tags TechnoMedic
// @Accept json
// @Produce json
// @Param no_order path string true "Order number"
// @Success 200 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 404 {object} externalEntity.TechnoMedicOrderResponse
// @Failure 500 {object} externalEntity.TechnoMedicOrderResponse
// @Router /technomedic/order/{no_order} [get]
func (h *TechnoMedicHandler) GetOrder(c echo.Context) error {
	noOrder := c.Param("no_order")
	if noOrder == "" {
		return c.JSON(http.StatusBadRequest, externalEntity.TechnoMedicOrderResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Order number is required",
		})
	}

	order, err := h.usecase.GetOrder(c.Request().Context(), noOrder)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "order not found: "+noOrder {
			statusCode = http.StatusNotFound
		}

		return c.JSON(statusCode, externalEntity.TechnoMedicOrderResponse{
			Code:    statusCode,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, externalEntity.TechnoMedicOrderResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Order retrieved successfully",
		Data:    order,
	})
}
