package delivery

import (
	"fmt"
	"slices"

	"github.com/oibacidem/lims-hl-seven/internal/delivery/serial/coax"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	analyxpanca "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_panca"
	analyxtrias "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_trias"
	swelabalfa "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_alfa"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type DeviceServerStrategy struct {
	coaxHandler        *coax.Handler
	defaultHandler     *tcp.HlSevenHandler
	analyxTriaHandler  *analyxtrias.Handler
	analyxPancaHandler *analyxpanca.Handler
	swelabAlfaHandler  *swelabalfa.Handler
}

func NewDeviceServerStrategy(
	coaxHandler *coax.Handler,
	defaultHandler *tcp.HlSevenHandler,
	analyxTriaHandler *analyxtrias.Handler,
	analyxPancaHandler *analyxpanca.Handler,
	swelabAlfaHandler *swelabalfa.Handler,
) *DeviceServerStrategy {
	return &DeviceServerStrategy{
		coaxHandler:        coaxHandler,
		defaultHandler:     defaultHandler,
		analyxTriaHandler:  analyxTriaHandler,
		analyxPancaHandler: analyxPancaHandler,
		swelabAlfaHandler:  swelabAlfaHandler,
	}
}

var _ repository.DeviceServerStrategy = (*DeviceServerStrategy)(nil)
var _ repository.DeviceSerialHandlerStrategy = (*DeviceServerStrategy)(nil)
var _ repository.DeviceTCPHandlerStrategy = (*DeviceServerStrategy)(nil)

// Make sure all device type is mapped
// this is only a compile check
// TODO: change this to unit test in the future
func init() {
	allDeviceType := append(serialDeviceType, tcpDeviceType...)
	for _, deviceType := range allDeviceType {
		d := DeviceServerStrategy{}
		_, err := d.ChooseDeviceServer(entity.Device{Type: deviceType})
		if err != nil {
			panic(fmt.Sprintf("device type %s is not supported", deviceType))
		}
	}
}

var serialDeviceType = []entity.DeviceType{
	entity.DeviceTypeCoax,
}

var tcpDeviceType = []entity.DeviceType{
	entity.DeviceTypeBA200,
	entity.DeviceTypeBA400,
	entity.DeviceTypeAnalyxTria,
	entity.DeviceTypeAnalyxPanca,
	entity.DeviceTypeSwelabAlfa,
	entity.DeviceTypeOther,
}

var deviceTypeNotSupport = []entity.DeviceType{
	entity.DeviceTypeA15,
}

func (d *DeviceServerStrategy) ChooseDeviceServer(device entity.Device) (server.Controller, error) {
	switch {
	case slices.Contains(deviceTypeNotSupport, device.Type):
		return nil, entity.ErrDeviceTypeNotSupport
	case slices.Contains(serialDeviceType, device.Type):
		s := server.NewSerial(device.ReceivePort, device.BaudRate)
		h, err := d.ChooseDeviceSerialHandler(device)
		if err != nil {
			return nil, err
		}
		s.SetHandler(h)
		return s, nil
	case slices.Contains(tcpDeviceType, device.Type):
		s := server.NewTCP(device.ReceivePort)
		h, err := d.ChooseDeviceTCPHandler(device)
		if err != nil {
			return nil, err
		}
		s.SetHandler(h)
		return s, nil
	default:
		return nil, entity.ErrDeviceTypeNotSupport
	}
}

func (d *DeviceServerStrategy) ChooseDeviceSerialHandler(device entity.Device) (server.SerialHandler, error) {
	// TODO: Add more device types here and change the default handler to the correct one
	switch device.Type {
	case entity.DeviceTypeCoax:
		return d.coaxHandler, nil
	default:
		return nil, entity.ErrDeviceTypeNotSupport
	}
}

func (d *DeviceServerStrategy) ChooseDeviceTCPHandler(device entity.Device) (server.TCPHandler, error) {
	// TODO: Add more device types here and change the default handler to the correct one
	switch device.Type {
	case entity.DeviceTypeBA200, entity.DeviceTypeBA400:
		return d.defaultHandler, nil
	case entity.DeviceTypeOther:
		return d.defaultHandler, nil
	case entity.DeviceTypeAnalyxTria:
		return d.analyxTriaHandler, nil
	case entity.DeviceTypeAnalyxPanca:
		return d.analyxPancaHandler, nil
	case entity.DeviceTypeSwelabAlfa:
		return d.swelabAlfaHandler, nil
	default:
		return d.defaultHandler, nil
	}
}
