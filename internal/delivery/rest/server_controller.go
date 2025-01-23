package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	configrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/config"
)

type (
	// ServerControllerHandler are handler to control the server and or listener
	ServerControllerHandler struct {
		Cfg       *configrepo.Repository
		TCPServer TCPServerController
	}

	// Some
	TCPServerController interface {
		SetPort(port string)
		Start() error
		State() string
		Serve()
		Stop() error
	}
)

func (s *ServerControllerHandler) RegisterRoute(p *echo.Group) {
	p.GET("/status", s.AllStatus)

	t := p.Group("/hl7tcp")
	t.GET("/status", s.TCPStatus)
	t.POST("/start", s.StartTCP)
	t.POST("/stop", s.StopTCP)
}

func (s *ServerControllerHandler) AllStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"rest":   "serving",
		"hl7tcp": s.TCPServer.State(),
	})
}

func (s *ServerControllerHandler) TCPStatus(c echo.Context) error {
	state := s.TCPServer.State()
	return c.JSON(http.StatusOK, state)
}

func (s *ServerControllerHandler) StartTCP(c echo.Context) error {
	ctx := c.Request().Context()

	port, err := s.Cfg.FindOne(ctx, "tcp_port")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	s.TCPServer.SetPort(port.Value)

	err = s.TCPServer.Start()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	go s.TCPServer.Serve()

	return c.JSON(http.StatusOK, "HL7 TCP server started")
}

func (s *ServerControllerHandler) StopTCP(c echo.Context) error {
	err := s.TCPServer.Stop()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "stopped")
}
