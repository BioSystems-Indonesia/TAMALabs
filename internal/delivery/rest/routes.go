package rest

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/BioSystems-Indonesia/TAMALabs/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	appMiddleware "github.com/BioSystems-Indonesia/TAMALabs/internal/middleware"
	"golang.org/x/exp/slices"
)

// Handler is a struct that contains the handler of the REST server.
type Handler struct {
	*HlSevenHandler
	*HealthCheckHandler
	*HealthHandler
	*PatientHandler
	*SpecimenHandler
	*WorkOrderHandler
	*FeatureListHandler
	*ObservationRequestHandler
	*TestTypeHandler
	*ResultHandler
	*ConfigHandler
	*UnitHandler
	*LogHandler
	*LicenseHandler
}

var blackListLoggingOnEndpoint = []string{
	// This is healthcheck endpoint, so we don't need to log it.
	"/api/v1/server/status",
}

// getCORSOrigins returns the allowed CORS origins
func getCORSOrigins() []string {
	// Default development origins
	defaultOrigins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:3000",
	}

	// Check for environment variable
	if corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); corsOrigins != "" {
		origins := strings.Split(corsOrigins, ",")
		// Trim whitespace from each origin
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
		return origins
	}

	return defaultOrigins
}

func RegisterMiddleware(e *echo.Echo) {
	slog.Info("registering middleware")

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogHost:    true,
		LogLatency: true,
		LogError:   true,
		LogMethod:  true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			if slices.Contains(blackListLoggingOnEndpoint, values.URI) {
				return nil
			}

			if values.Error != nil {
				slog.Error("request error",
					"method", values.Method,
					"uri", values.URI,
					"status", values.Status,
					"latency", values.Latency,
					"error", values.Error,
				)
			} else if values.Status >= http.StatusBadRequest {
				slog.Error("request error",
					"method", values.Method,
					"uri", values.URI,
					"status", values.Status,
					"latency", values.Latency,
				)
			} else {
				slog.Debug("request",
					"method", values.Method,
					"uri", values.URI,
					"status", values.Status,
					"latency", values.Latency,
					"error", values.Error,
				)
			}

			return nil
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     getCORSOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		ExposeHeaders:    []string{"X-Total-Count", "Content-Range", "Accept-Ranges", "Content-Length", "Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
}

// RegisterRoutes registers the routes of the REST server.
func RegisterRoutes(
	e *echo.Echo,
	handler *Handler,
	deviceHandler *DeviceHandler,
	serverControllerHandler *ServerControllerHandler,
	testTemplateHandler *TestTemplateHandler,
	adminHandler *AdminHandler,
	authHandler *AuthHandler,
	roleHandler *RoleHandler,
	khanzaExternalHandler *KhanzaExternalHandler,
	externalHandler *ExternalHandler,
	authMiddleware *appMiddleware.JWTMiddleware,
) {
	slog.Info("registering routes")

	registerFrontendPath(e)

	api := e.Group("/api")
	unauthenticatedV1 := api.Group("/v1")
	unauthenticatedV1.GET("/ping", handler.Ping)
	unauthenticatedV1.POST("/login", authHandler.Login)
	unauthenticatedV1.POST("/logout", authHandler.Logout)

	// Add health endpoint (unauthenticated for monitoring)
	handler.HealthHandler.RegisterRoutes(unauthenticatedV1)

	// Add license endpoint (unauthenticated for license checking)
	handler.LicenseHandler.RegisterRoutes(unauthenticatedV1)

	authenticatedV1 := api.Group("/v1", authMiddleware.Middleware())
	authenticatedV1.GET("/check-auth", handler.Ping)
	authenticatedV1.GET("/profile", authHandler.GetProfile)
	authenticatedV1.GET("/permissions", authHandler.GetPermissions)
	authenticatedV1.POST("/logout", authHandler.Logout)
	patient := authenticatedV1.Group("/patient")
	{
		patient.GET("", handler.FindPatients)
		patient.GET("/:id", handler.GetOnePatient)
		patient.GET("/:id/result/history", handler.GetPatientResultHistory)
		patient.POST("", handler.CreatePatient)
		patient.PUT("/:id", handler.UpdatePatient)
		patient.DELETE("/:id", handler.DeletePatient)
	}

	specimen := authenticatedV1.Group("/specimen")
	{
		specimen.GET("", handler.FindSpecimens)
		specimen.GET("/:id", handler.GetOneSpecimen)
	}

	observationRequest := authenticatedV1.Group("/observation-request")
	{
		observationRequest.GET("", handler.FindObservationRequests)
		observationRequest.GET("/:id", handler.GetOneObservationRequest)
	}

	workOrder := authenticatedV1.Group("/work-order")
	{
		workOrder.GET("", handler.FindWorkOrders)
		workOrder.GET("/barcode", handler.GetWorkOrderBarcode)
		workOrder.POST("", handler.CreateWorkOrder)
		workOrder.POST("/run", handler.RunWorkOrder)
		workOrder.POST("/cancel", handler.CancelOrder)
		workOrder.GET("/:id", handler.GetOneWorkOrder)
		workOrder.PUT("/:id", handler.EditWorkOrder)
		workOrder.DELETE("/:id", handler.DeleteWorkOrder)
	}

	deviceHandler.RegisterRoute(authenticatedV1.Group("/device"))

	serverControllerHandler.RegisterRoute(authenticatedV1.Group("/server"))

	testType := authenticatedV1.Group("/test-type")
	{
		testType.GET("", handler.ListTestType)
		testType.GET("/filter", handler.ListTestTypeFilter)
		testType.GET("/:id", handler.GetOneTestType)
		testType.GET("/code/:code", handler.GetOneTestTypeByCode)
		testType.GET("/alias-code/:alias_code", handler.GetOneTestTypeByAliasCode)
		testType.POST("/upload", handler.UploadBulkTestType)
		testType.POST("", handler.CreateTestType)
		testType.PUT("/:id", handler.UpdateTestType)
		testType.DELETE("/:id", handler.DeleteTestType)
	}

	testTemplate := authenticatedV1.Group("/test-template")
	{
		testTemplate.GET("", testTemplateHandler.ListTestTemplate)
		testTemplate.GET("/:id", testTemplateHandler.GetOneTestTemplate)
		testTemplate.POST("", testTemplateHandler.CreateTestTemplate)
		testTemplate.PUT("/:id", testTemplateHandler.UpdateTestTemplate)
		testTemplate.PUT("/:id/update-diff", testTemplateHandler.CheckUpdateDifferenceTestTemplate)
		testTemplate.DELETE("/:id", testTemplateHandler.DeleteTestTemplate)
	}

	result := authenticatedV1.Group("/result")
	{
		result.GET("", handler.ListResult)
		result.POST("/refresh", handler.RefreshResult)
		result.GET("/:work_order_id", handler.GetResult)
		result.POST("/:work_order_id/approve", handler.ApproveResult)
		result.POST("/:work_order_id/reject", handler.RejectResult)
		result.PUT("/:work_order_id/test", handler.AddTestResult)
		result.PUT("/:work_order_id/test/:test_result_id/pick", handler.TooglePickTestResult)
		result.DELETE("/:work_order_id/test/:test_result_id", handler.DeleteTestResult)
		result.POST("/egfr/calculate", handler.CalculateEGFR)
	}
	resultUnauthenticated := unauthenticatedV1.Group("/result")
	{
		resultUnauthenticated.POST("/a15/upload", handler.UploadFileA15)
	}

	config := authenticatedV1.Group("/config")
	{
		config.GET("", handler.ListConfig)
		config.GET("/:key", handler.GetConfig)
		config.PUT("/:key", handler.EditConfig)
	}

	unit := authenticatedV1.Group("/unit")
	{
		unit.GET("", handler.ListUnit)
	}

	admin := authenticatedV1.Group("/user")
	{
		admin.GET("", adminHandler.FindAdmins)
		admin.GET("/:id", adminHandler.GetOneAdmin)
		admin.POST("", adminHandler.CreateAdmin)
		admin.PUT("/:id", adminHandler.UpdateAdmin)
		admin.DELETE("/:id", adminHandler.DeleteAdmin)
	}

	role := authenticatedV1.Group("/role")
	{
		role.GET("", roleHandler.FindRoles)
		role.GET("/:id", roleHandler.GetOneRole)
	}

	unauthenticatedLog := unauthenticatedV1.Group("/log")
	{
		unauthenticatedLog.GET("/stream", handler.LogHandler.StreamLog)
	}

	log := authenticatedV1.Group("/log")
	{
		log.GET("/export", handler.LogHandler.ExportLog)
	}

	khanzaExternalHandler.RegisterRoutes(unauthenticatedV1)
	handler.RegisterFeatureList(authenticatedV1)
	externalHandler.RegisterRoutes(authenticatedV1)
}

func registerFrontendPath(e *echo.Echo) {
	h := http.FileServer(http.FS(web.Content()))
	e.GET("/*", echo.WrapHandler(h))
}
