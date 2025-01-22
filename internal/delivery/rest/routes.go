package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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
}

func RegisterMiddleware(e *echo.Echo) {
	log.Info("Registering middleware")

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogHost:    true,
		LogLatency: true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {

			json := log.JSON{
				"method":  values.Method,
				"uri":     values.URI,
				"status":  values.Status,
				"latency": values.Latency,
			}
			if values.Error != nil {
				json["error"] = values.Error.Error()
			}

			log.Infoj(json)

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
func RegisterRoutes(e *echo.Echo, handler *Handler, deviceHandler *DeviceHandler) {
	log.Info("Registering routes")

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
		workOrder.POST("/:id/show/add-test", handler.AddTestWorkOrder)
		workOrder.GET("/:id", handler.GetOneWorkOrder)
		workOrder.DELETE("/:id/test/:patient_id", handler.DeleteTestWorkOrder)
		workOrder.DELETE("/:id", handler.DeleteWorkOrder)
	}

	device := v1.Group("/device")
	{
		device.GET("", deviceHandler.ListDevices)
		device.POST("", deviceHandler.CreateDevice)
		device.GET("/:id", deviceHandler.GetDevice)
		device.PUT("/:id", deviceHandler.UpdateDevice)
		device.DELETE("/:id", deviceHandler.DeleteDevice)
	}

	testType := v1.Group("/test-type")
	{
		testType.GET("", handler.ListTestType)
		testType.GET("/:id", handler.GetOneTestType)
		testType.POST("", handler.CreateTestType)
		testType.PUT("/:id", handler.UpdateTestType)
	}

	result := v1.Group("/result")
	{
		result.GET("", handler.ListResult)
		result.GET("/:barcode", handler.GetResult)
	}

	config := v1.Group("/config")
	{
		config.GET("", handler.ListConfig)
		config.GET("/:key", handler.GetConfig)
		config.PUT("/:key", handler.EditConfig)
	}

	handler.RegisterFeatureList(v1)
}

func registerFrontendPath(e *echo.Echo) {
	h := http.FileServer(http.FS(web.Content()))
	e.GET("/*", echo.WrapHandler(h))
}
