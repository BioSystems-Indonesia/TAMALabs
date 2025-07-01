package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	configrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/config"
	tcpserver "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/server"
	deviceuc "github.com/oibacidem/lims-hl-seven/internal/usecase/device"
)

type (
	// ServerControllerHandler are handler to control the server and or listener
	ServerControllerHandler struct {
		cfg           *configrepo.Repository
		tcpserver     *tcpserver.TCPServer
		deviceUsecase *deviceuc.DeviceUseCase
	}
)

func NewServerControllerHandler(
	cfg *configrepo.Repository,
	tcpserver *tcpserver.TCPServer,
) *ServerControllerHandler {
	return &ServerControllerHandler{
		cfg:       cfg,
		tcpserver: tcpserver,
	}
}

func (s *ServerControllerHandler) RegisterRoute(p *echo.Group) {
	p.GET("/status", s.AllStatus)

	t := p.Group("/hl7tcp")
	t.GET("/status", s.TCPStatus)
	t.POST("/start/:device_id", s.StartTCP)
	t.POST("/stop/:device_id", s.StopTCP)
}

func (s *ServerControllerHandler) AllStatus(c echo.Context) error {
	mapState := s.tcpserver.GetAllServerState()

	return c.JSON(http.StatusOK, map[string]any{
		"rest":   constant.ServerStateServing,
		"hl7tcp": mapState,
	})
}

func (s *ServerControllerHandler) TCPStatus(c echo.Context) error {
	mapState := s.tcpserver.GetAllServerState()

	return c.JSON(http.StatusOK, mapState)
}

func (s *ServerControllerHandler) StartTCP(c echo.Context) error {
	ctx := c.Request().Context()
	deviceID, err := strconv.Atoi(c.Param("device_id"))

	device, err := s.deviceUsecase.FindOneByID(ctx, int64(deviceID))
	if err != nil {
		return handleError(c, err)
	}

	_, err = s.tcpserver.StartNewServer(ctx, device)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, "HL7 TCP server started")
}

func (s *ServerControllerHandler) StopTCP(c echo.Context) error {
	ctx := c.Request().Context()
	deviceID, err := strconv.Atoi(c.Param("device_id"))

	err = s.tcpserver.StopServerByDeviceID(ctx, deviceID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, "stopped")
}
