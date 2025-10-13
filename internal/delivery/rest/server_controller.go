package rest

import (
	"net/http"
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	serverrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/server"
	configrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/config"
	deviceuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/device"
	"github.com/labstack/echo/v4"
)

type (
	// ServerControllerHandler are handler to control the server and or listener
	ServerControllerHandler struct {
		cfg           *configrepo.Repository
		serverRepo    *serverrepo.ControllerRepository
		deviceUsecase *deviceuc.DeviceUseCase
	}
)

func NewServerControllerHandler(
	cfg *configrepo.Repository,
	serverRepo *serverrepo.ControllerRepository,
) *ServerControllerHandler {
	return &ServerControllerHandler{
		cfg:        cfg,
		serverRepo: serverRepo,
	}
}

func (s *ServerControllerHandler) RegisterRoute(p *echo.Group) {
	p.GET("/status", s.AllStatus)
	p.GET("/serial-port-list", s.GetSerialPortList)

	t := p.Group("/hl7tcp")
	t.GET("/status", s.TCPStatus)
	t.POST("/start/:device_id", s.StartTCP)
	t.POST("/stop/:device_id", s.StopTCP)
}

func (s *ServerControllerHandler) AllStatus(c echo.Context) error {
	mapState := s.serverRepo.GetAllServerState()

	return c.JSON(http.StatusOK, map[string]any{
		"rest":   constant.ServerStateServing,
		"hl7tcp": mapState,
	})
}

func (s *ServerControllerHandler) TCPStatus(c echo.Context) error {
	mapState := s.serverRepo.GetAllServerState()

	return c.JSON(http.StatusOK, mapState)
}

func (s *ServerControllerHandler) StartTCP(c echo.Context) error {
	ctx := c.Request().Context()
	deviceID, _ := strconv.Atoi(c.Param("device_id"))

	device, err := s.deviceUsecase.FindOneByID(ctx, int64(deviceID))
	if err != nil {
		return handleError(c, err)
	}

	_, err = s.serverRepo.StartNewServer(ctx, device)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, "HL7 TCP server started")
}

func (s *ServerControllerHandler) StopTCP(c echo.Context) error {
	ctx := c.Request().Context()
	deviceID, _ := strconv.Atoi(c.Param("device_id"))

	err := s.serverRepo.StopServerByDeviceID(ctx, deviceID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, "stopped")
}

func (s *ServerControllerHandler) GetSerialPortList(c echo.Context) error {
	portList, err := s.serverRepo.GetAllSerialPorts()
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, entity.PaginationResponse[entity.Table]{
		Data:  mapToTable(portList),
		Total: int64(len(portList)),
	})
}

func mapToTable(portList []string) []entity.Table {
	table := make([]entity.Table, len(portList))
	for i, port := range portList {
		table[i] = entity.Table{
			ID:   port,
			Name: port,
		}
	}
	return table
}
