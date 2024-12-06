package rest

import "github.com/labstack/echo/v4"

// Handler is a struct that contains the handler of the REST server.
type Handler struct {
	*HlSevenHandler
}

// RegisterRoutes registers the routes of the REST server.
func RegisterRoutes(e *echo.Echo, handler *Handler) *echo.Echo {
	v1 := e.Group("/v1")

	// HL Seven routes
	hlSeven := v1.Group("/hl-seven")
	// Register the routes here
	hlSeven.POST("/orm", handler.SendORM)

	return e
}
