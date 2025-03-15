package rest

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oibacidem/lims-hl-seven/web"
)

// Handler is a struct that contains the handler of the REST server.
type Handler struct {
	*HlSevenHandler
	*HealthCheckHandler
	*PatientHandler
	*SpecimenHandler
	*WorkOrderHandler
	*FeatureListHandler
	*ObservationRequestHandler
	*TestTypeHandler
	*ResultHandler
	*ConfigHandler
	*UnitHandler
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
			slog.Info("request",
				"method", values.Method,
				"uri", values.URI,
				"status", values.Status,
				"latency", values.Latency,
				"error", values.Error,
			)

			return nil
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
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
) {
	slog.Info("registering routes")

	registerFrontendPath(e)

	api := e.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/ping", handler.Ping)

	patient := v1.Group("/patient")
	{
		patient.GET("", handler.FindPatients)
		patient.GET("/:id", handler.GetOnePatient)
		patient.POST("", handler.CreatePatient)
		patient.PUT("/:id", handler.UpdatePatient)
		patient.DELETE("/:id", handler.DeletePatient)
	}

	specimen := v1.Group("/specimen")
	{
		specimen.GET("", handler.FindSpecimens)
		specimen.GET("/:id", handler.GetOneSpecimen)
	}

	observationRequest := v1.Group("/observation-request")
	{
		observationRequest.GET("", handler.FindObservationRequests)
		observationRequest.GET("/:id", handler.GetOneObservationRequest)
	}

	workOrder := v1.Group("/work-order")
	{
		workOrder.GET("", handler.FindWorkOrders)
		workOrder.POST("", handler.CreateWorkOrder)
		workOrder.POST("/run", handler.RunWorkOrder)
		workOrder.POST("/cancel", handler.CancelOrder)
		workOrder.GET("/:id", handler.GetOneWorkOrder)
		workOrder.PUT("/:id", handler.EditWorkOrder)
		workOrder.DELETE("/:id", handler.DeleteWorkOrder)
	}

	deviceHandler.RegisterRoute(v1.Group("/device"))

	serverControllerHandler.RegisterRoute(v1.Group("/server"))

	testType := v1.Group("/test-type")
	{
		testType.GET("", handler.ListTestType)
		testType.GET("/filter", handler.ListTestTypeFilter)
		testType.GET("/:id", handler.GetOneTestType)
		testType.POST("", handler.CreateTestType)
		testType.PUT("/:id", handler.UpdateTestType)
		testType.DELETE("/:id", handler.DeleteTestType)
	}

	testTemplate := v1.Group("/test-template")
	{
		testTemplate.GET("", testTemplateHandler.ListTestTemplate)
		testTemplate.GET("/:id", testTemplateHandler.GetOneTestTemplate)
		testTemplate.POST("", testTemplateHandler.CreateTestTemplate)
		testTemplate.PUT("/:id", testTemplateHandler.UpdateTestTemplate)
		testTemplate.DELETE("/:id", testTemplateHandler.DeleteTestTemplate)
	}

	result := v1.Group("/result")
	{
		result.GET("", handler.ListResult)
		result.GET("/:work_order_id", handler.GetResult)
		result.PUT("/:work_order_id/test", handler.AddTestResult)
		result.PUT("/:work_order_id/test/:test_result_id/pick", handler.TooglePickTestResult)
		result.DELETE("/:work_order_id/test/:test_result_id", handler.DeleteTestResult)
	}

	config := v1.Group("/config")
	{
		config.GET("", handler.ListConfig)
		config.GET("/:key", handler.GetConfig)
		config.PUT("/:key", handler.EditConfig)
	}

	unit := v1.Group("/unit")
	{
		unit.GET("", handler.ListUnit)
	}

	handler.RegisterFeatureList(v1)
}

func registerFrontendPath(e *echo.Echo) {
	h := http.FileServer(http.FS(web.Content()))
	e.GET("/*", echo.WrapHandler(h))
}
