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
	// Documentation endpoint - publicly accessible, no middleware
	router.GET("/technomedic/documentation", h.GetDocumentation)

	// API endpoints - protected by integration check middleware
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

// GetDocumentation serves Swagger UI for TechnoMedic API documentation
func (h *TechnoMedicHandler) GetDocumentation(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TechnoMedic API Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
        .topbar {
            display: none;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const spec = ` + "`" + `
openapi: 3.0.3
info:
  title: TechnoMedic API Integration
  description: |
    API bridging untuk integrasi dengan TechnoMedic SIMRS.
    
    ## Configuration
    Sebelum menggunakan API ini, pastikan TechnoMedic integration sudah diaktifkan:
    1. Login sebagai Administrator
    2. Navigate to Config page
    3. Enable "SIMRS Bridging"
    4. Select "TechnoMedic (API)"
    5. Save configuration
    
    ## Features
    - Get master data (test types, sub-categories, doctors, analysts)
    - Create lab orders
    - Retrieve order results
    - Support multiple test selection methods (IDs, codes, sub-categories)
    
  version: 1.0.0
  contact:
    name: API Support
  license:
    name: Proprietary
    
servers:
  - url: ` + "`" + ` + window.location.origin + ` + "`" + `/api/v1
    description: Current Server

tags:
  - name: Master Data
    description: Endpoints untuk mendapatkan master data
  - name: Orders
    description: Endpoints untuk manajemen order laboratorium

paths:
  /technomedic/test-types:
    get:
      tags:
        - Master Data
      summary: Get All Test Types
      description: Mendapatkan daftar semua test types yang tersedia
      operationId: getTestTypes
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  code:
                    type: integer
                    example: 200
                  status:
                    type: string
                    example: success
                  message:
                    type: string
                    example: Test types retrieved successfully
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/TestType'

  /technomedic/sub-categories:
    get:
      tags:
        - Master Data
      summary: Get All Sub-Categories
      description: Mendapatkan daftar semua sub-categories dari tabel master
      operationId: getSubCategories
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  code:
                    type: integer
                    example: 200
                  status:
                    type: string
                    example: success
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/SubCategoryInfo'

  /technomedic/sub-categories/{id}/test-types:
    get:
      tags:
        - Master Data
      summary: Get Test Types by Sub-Category
      description: Mendapatkan daftar test types untuk sub-category tertentu
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Successful response

  /technomedic/doctors:
    get:
      tags:
        - Master Data
      summary: Get All Doctors
      description: Mendapatkan daftar semua dokter yang aktif
      responses:
        '200':
          description: Successful response

  /technomedic/analysts:
    get:
      tags:
        - Master Data
      summary: Get All Analysts
      description: Mendapatkan daftar semua analis yang aktif
      responses:
        '200':
          description: Successful response

  /technomedic/order:
    post:
      tags:
        - Orders
      summary: Create Lab Order
      description: |
        Membuat order laboratorium baru dari TechnoMedic.
        
        ## Test Selection Methods
        Minimal salah satu dari method berikut harus diisi:
        - test_type_ids: Array of test type IDs (RECOMMENDED)
        - sub_category_ids: Array of sub-category IDs (RECOMMENDED)
        - param_request: Array of test type codes
        - sub_category_request: Array of sub-category names
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateOrderRequest'
            examples:
              withTestTypeIDs:
                summary: Using test_type_ids
                value:
                  no_order: TM-2024-001
                  patient:
                    patient_id: P001
                    full_name: John Doe
                    sex: M
                    birthdate: '1990-01-15'
                    medical_record_number: MR001
                  test_type_ids: [1, 2, 3]
      responses:
        '201':
          description: Order created successfully

  /technomedic/order/{no_order}:
    get:
      tags:
        - Orders
      summary: Get Order Details
      description: Mendapatkan detail order termasuk hasil pemeriksaan
      parameters:
        - name: no_order
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Successful response

components:
  schemas:
    TestType:
      type: object
      properties:
        id:
          type: string
          example: '1'
        code:
          type: string
          example: HB
        name:
          type: string
          example: Hemoglobin
        category:
          type: string
          example: Hematologi
        sub_category:
          type: string
          example: Complete Blood Count
        specimen_type:
          type: string
          example: Whole Blood
        unit:
          type: string
          example: g/dL

    SubCategoryInfo:
      type: object
      properties:
        id:
          type: string
          example: '1'
        name:
          type: string
          example: Complete Blood Count
        category:
          type: string
          example: Hematologi
        description:
          type: string
          example: Pemeriksaan darah lengkap

    CreateOrderRequest:
      type: object
      required:
        - no_order
        - patient
      properties:
        no_order:
          type: string
          example: TM-2024-001
        patient:
          type: object
          required:
            - patient_id
            - full_name
            - sex
            - birthdate
          properties:
            patient_id:
              type: string
              example: P001
            full_name:
              type: string
              example: John Doe
            sex:
              type: string
              enum: [M, F]
              example: M
            birthdate:
              type: string
              format: date
              example: '1990-01-15'
            medical_record_number:
              type: string
              example: MR001
            phone_number:
              type: string
              example: '081234567890'
        test_type_ids:
          type: array
          items:
            type: integer
          example: [1, 2, 3]
        sub_category_ids:
          type: array
          items:
            type: integer
          example: [1, 2]
        param_request:
          type: array
          items:
            type: string
          example: [HB, WBC]
` + "`" + `;

            SwaggerUIBundle({
                spec: jsyaml.load(spec),
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                tryItOutEnabled: true,
                supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch']
            });
        };
    </script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/js-yaml/4.1.0/js-yaml.min.js"></script>
</body>
</html>`

	return c.HTML(http.StatusOK, html)
}
