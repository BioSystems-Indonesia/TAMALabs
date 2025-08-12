package delivery

import (
	"fmt"
	"slices"

	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/a15"
	swelabalfa "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_alfa"
	swelablumi "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_lumi"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type DeviceServerStrategy struct {
	a15Handler *a15.Handler
	// coaxHandler        *coax.Handler
	// ncc3300            *ncc3300.Handler
	defaultHandler *tcp.HlSevenHandler
	// analyxTriaHandler  *analyxtrias.Handler
	// analyxPancaHandler *analyxpanca.Handler
	swelabAlfaHandler *swelabalfa.Handler
	swelabAlfaBasic   *swelabalfa.Handler
	// swelabLumiHandler  *swelablumi.Handler
	// alifaxHandler      *alifax.Handler
	// ncc61Handler       *ncc61.Handler
}

func NewDeviceServerStrategy(
	a15Handler *a15.Handler,
	// coaxHandler *coax.Handler,
	// ncc3300 *ncc3300.Handler,
	defaultHandler *tcp.HlSevenHandler,
	// analyxTriaHandler *analyxtrias.Handler,
	// analyxPancaHandler *analyxpanca.Handler,
	swelabAlfaHandler *swelabalfa.Handler,
	swelabLumiHandler *swelablumi.Handler,
	// alifaxHandler *alifax.Handler,
	// ncc61handler *ncc61.Handler,
) *DeviceServerStrategy {
	return &DeviceServerStrategy{
		a15Handler: a15Handler,
		// coaxHandler:        coaxHandler,
		// ncc3300:            ncc3300,
		defaultHandler: defaultHandler,
		// analyxTriaHandler:  analyxTriaHandler,
		// analyxPancaHandler: analyxPancaHandler,
		swelabAlfaHandler: swelabAlfaHandler,
		swelabAlfaBasic:   swelabAlfaHandler,
		// swelabLumiHandler:  swelabLumiHandler,
		// alifaxHandler:      alifaxHandler,
		// ncc61Handler:       ncc61handler,
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
	// entity.DeviceTypeCoax,
	// entity.DeviceTypeBiomedicaNCC3300,
	// entity.DeviceTypeAlifax,
}

var tcpDeviceType = []entity.DeviceType{
	entity.DeviceTypeA15,
	// entity.DeviceTypeBA200,
	// entity.DeviceTypeBA400,
	// entity.DeviceTypeAnalyxTria,
	// entity.DeviceTypeAnalyxPanca,
	entity.DeviceTypeSwelabAlfa,
	entity.DeviceTypeSwelabBasic,
	// entity.DeviceTypeSwelabLumi,
	// entity.DeviceTypeBiomedicaNCC61,
	entity.DeviceTypeOther,
}

var deviceTypeNotSupport = []entity.DeviceType{
	// entity.DeviceTypeA15,
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
	// case entity.DeviceTypeCoax:
	// 	return d.coaxHandler, nil
	// case entity.DeviceTypeBiomedicaNCC3300:
	// 	return d.ncc3300, nil
	// case entity.DeviceTypeAlifax:
	// 	return d.alifaxHandler, nil

	default:
		return nil, entity.ErrDeviceTypeNotSupport
	}
}

func (d *DeviceServerStrategy) ChooseDeviceTCPHandler(device entity.Device) (server.TCPHandler, error) {
	// TODO: Add more device types here and change the default handler to the correct one
	switch device.Type {
	// case entity.DeviceTypeBA200, entity.DeviceTypeBA400:
	// 	return d.defaultHandler, nil
	// case entity.DeviceTypeOther:
	// 	return d.defaultHandler, nil
	// case entity.DeviceTypeAnalyxTria:
	// 	return d.analyxTriaHandler, nil
	// case entity.DeviceTypeAnalyxPanca:
	// 	return d.analyxPancaHandler, nil
	case entity.DeviceTypeSwelabAlfa:
		return d.swelabAlfaHandler, nil
	case entity.DeviceTypeSwelabBasic:
		return d.swelabAlfaBasic, nil
	// case entity.DeviceTypeSwelabLumi:
	// 	return d.swelabLumiHandler, nil
	case entity.DeviceTypeA15:
		return d.a15Handler, nil
	// case entity.DeviceTypeBiomedicaNCC61:
	// 	return d.ncc61Handler, nil
	default:
		return d.defaultHandler, nil
	}
}
