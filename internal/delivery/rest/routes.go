package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	_ "github.com/oibacidem/lims-hl-seven/statik"
	"github.com/rakyll/statik/fs"
)

// Handler is a struct that contains the handler of the REST server.
type Handler struct {
	*HlSevenHandler
	*HealthCheckHandler
}

// RegisterRoutes registers the routes of the REST server.
func RegisterRoutes(e *echo.Echo, handler *Handler, cfg *config.Schema) *echo.Echo {
	registerFrontendPath(e)

	api := e.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/ping", handler.Ping)
	// HL Seven routes
	hlSeven := v1.Group("/hl-seven")
	// Register the routes here
	hlSeven.POST("/orm", handler.SendORM)

	return e
}

func registerFrontendPath(e *echo.Echo) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	h := http.FileServer(statikFS)
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", h)))
}
