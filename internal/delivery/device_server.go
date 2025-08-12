package delivery

import (
	"fmt"
	"slices"

	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/a15"
	swelabalfa "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_alfa"
	swelablumi "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_lumi"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type DeviceServerStrategy struct {
	a15Handler        *a15.Handler
	swelabAlfaHandler *swelabalfa.Handler
	swelabAlfaBasic   *swelabalfa.Handler
}

func NewDeviceServerStrategy(
	a15Handler *a15.Handler,
	swelabAlfaHandler *swelabalfa.Handler,
	swelabLumiHandler *swelablumi.Handler,
) *DeviceServerStrategy {
	return &DeviceServerStrategy{
		a15Handler:        a15Handler,
		swelabAlfaHandler: swelabAlfaHandler,
		swelabAlfaBasic:   swelabAlfaHandler,
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

var serialDeviceType = []entity.DeviceType{}

var tcpDeviceType = []entity.DeviceType{
	entity.DeviceTypeA15,
	entity.DeviceTypeSwelabAlfa,
	entity.DeviceTypeSwelabBasic,
}

var deviceTypeNotSupport = []entity.DeviceType{}

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
	switch device.Type {
	default:
		return nil, entity.ErrDeviceTypeNotSupport
	}
}

func (d *DeviceServerStrategy) ChooseDeviceTCPHandler(device entity.Device) (server.TCPHandler, error) {
	// TODO: Add more device types here and change the default handler to the correct one
	switch device.Type {
	case entity.DeviceTypeSwelabAlfa:
		return d.swelabAlfaHandler, nil
	case entity.DeviceTypeSwelabBasic:
		return d.swelabAlfaBasic, nil
	case entity.DeviceTypeA15:
		return d.a15Handler, nil
	default:
		return nil, nil
	}
}
