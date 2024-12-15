package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/oibacidem/lims-hl-seven/statik"
	"github.com/rakyll/statik/fs"
)

// Handler is a struct that contains the handler of the REST server.
type Handler struct {
	*HlSevenHandler
	*HealthCheckHandler
	*PatientHandler
	*SpecimentHandler
	*WorkOrderHandler
	*FeatureListHandler
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
func RegisterRoutes(e *echo.Echo, handler *Handler) {
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

	speciment := v1.Group("/speciment")
	{
		speciment.GET("", handler.FindSpeciments)
		speciment.GET("/:id", handler.GetOneSpeciment)
		speciment.POST("", handler.CreateSpeciment)
		speciment.PUT("/:id", handler.UpdateSpeciment)
		speciment.DELETE("/:id", handler.DeleteSpeciment)
	}

	workOrder := v1.Group("/work-order")
	{
		workOrder.GET("", handler.FindWorkOrders)
		workOrder.POST("", handler.CreateWorkOrder)
		workOrder.POST("/:id/speciment", handler.AddSpeciment)
		workOrder.GET("/:id", handler.GetOneWorkOrder)
		workOrder.PUT("/:id", handler.UpdateWorkOrder)
		workOrder.DELETE("/:id", handler.DeleteWorkOrder)
	}

		handler.RegisterFeatureList(v1)
}

func registerFrontendPath(e *echo.Echo) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	h := http.FileServer(statikFS)
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", h)))
}
