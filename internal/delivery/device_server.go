package delivery

import (
	"fmt"
	"slices"

	"github.com/oibacidem/lims-hl-seven/internal/delivery/serial/alifax"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/serial/cbs400"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/serial/coax"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/serial/diestro"
	ncc3300 "github.com/oibacidem/lims-hl-seven/internal/delivery/serial/ncc_3300"
	verifyu120 "github.com/oibacidem/lims-hl-seven/internal/delivery/serial/verifyU120"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/a15"
	analyxpanca "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_panca"
	analyxtrias "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_trias"
	ncc61 "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/neomedika_ncc61"
	swelabalfa "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_alfa"
	swelablumi "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_lumi"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/wondfo"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type DeviceServerStrategy struct {
	a15Handler         *a15.Handler
	coaxHandler        *coax.Handler
	diestroHandler     *diestro.Handler
	ncc3300            *ncc3300.Handler
	defaultHandler     *tcp.HlSevenHandler
	analyxTriaHandler  *analyxtrias.Handler
	analyxPancaHandler *analyxpanca.Handler
	swelabAlfaHandler  *swelabalfa.Handler
	swelabAlfaBasic    *swelabalfa.Handler
	swelabLumiHandler  *swelablumi.Handler
	alifaxHandler      *alifax.Handler
	ncc61Handler       *ncc61.Handler

	wondfoHandler     *wondfo.Handler
	cbs400Handler     *cbs400.Handler
	verifyu120Handler *verifyu120.Handler
}

func NewDeviceServerStrategy(
	a15Handler *a15.Handler,
	coaxHandler *coax.Handler,
	diestroHandler *diestro.Handler,
	ncc3300 *ncc3300.Handler,
	defaultHandler *tcp.HlSevenHandler,
	analyxTriaHandler *analyxtrias.Handler,
	analyxPancaHandler *analyxpanca.Handler,
	swelabAlfaHandler *swelabalfa.Handler,
	swelabLumiHandler *swelablumi.Handler,
	alifaxHandler *alifax.Handler,
	ncc61handler *ncc61.Handler,

	wondfoHandler *wondfo.Handler,
	cbs400Handler *cbs400.Handler,
	verifyu120Handler *verifyu120.Handler,

) *DeviceServerStrategy {
	return &DeviceServerStrategy{
		a15Handler:         a15Handler,
		coaxHandler:        coaxHandler,
		diestroHandler:     diestroHandler,
		ncc3300:            ncc3300,
		defaultHandler:     defaultHandler,
		analyxTriaHandler:  analyxTriaHandler,
		analyxPancaHandler: analyxPancaHandler,
		swelabAlfaHandler:  swelabAlfaHandler,
		swelabAlfaBasic:    swelabAlfaHandler,
		swelabLumiHandler:  swelabLumiHandler,
		alifaxHandler:      alifaxHandler,
		ncc61Handler:       ncc61handler,

		wondfoHandler:     wondfoHandler,
		cbs400Handler:     cbs400Handler,
		verifyu120Handler: verifyu120Handler,
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
	entity.DeviceTypeDiestro,
	entity.DeviceTypeNeomedicaNCC3300,
	entity.DeviceTypeAlifax,
	entity.DeviceTypeCBS400,
	entity.DeviceTypeVerifyU120,
}

var tcpDeviceType = []entity.DeviceType{
	entity.DeviceTypeA15,
	entity.DeviceTypeBA200,
	entity.DeviceTypeBA400,
	entity.DeviceTypeAnalyxTria,
	entity.DeviceTypeAnalyxPanca,
	entity.DeviceTypeSwelabAlfa,
	entity.DeviceTypeSwelabBasic,
	entity.DeviceTypeSwelabLumi,
	entity.DeviceTypeNeomedicaNCC61,
	entity.DeviceTypeOther,
	entity.DeviceTypeWondfo,
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
	case entity.DeviceTypeDiestro:
		return d.diestroHandler, nil
	case entity.DeviceTypeCoax:
		return d.coaxHandler, nil
	case entity.DeviceTypeNeomedicaNCC3300:
		return d.ncc3300, nil
	case entity.DeviceTypeAlifax:
		return d.alifaxHandler, nil
	case entity.DeviceTypeVerifyU120:
		return d.verifyu120Handler, nil
	case entity.DeviceTypeCBS400:
		return d.cbs400Handler, nil

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
	case entity.DeviceTypeSwelabBasic:
		return d.swelabAlfaBasic, nil
	case entity.DeviceTypeSwelabLumi:
		return d.swelabLumiHandler, nil
	case entity.DeviceTypeA15:
		return d.a15Handler, nil
	case entity.DeviceTypeNeomedicaNCC61:
		return d.ncc61Handler, nil
	case entity.DeviceTypeWondfo:
		return d.wondfoHandler, nil
	default:
		return nil, nil
	}
}
