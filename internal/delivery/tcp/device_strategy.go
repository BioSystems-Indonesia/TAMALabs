package tcp

import (
	analyxpanca "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_panca"
	analyxtrias "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_trias"
	swelabalfa "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_alfa"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type DeviceStrategy struct {
	defaultHandler     *HlSevenHandler
	analyxTriaHandler  *analyxtrias.Handler
	analyxPancaHandler *analyxpanca.Handler
	swelabAlfaHandler  *swelabalfa.Handler
}

func NewDeviceStrategy(
	defaultHandler *HlSevenHandler,
	analyxTriaHandler *analyxtrias.Handler,
	analyxPancaHandler *analyxpanca.Handler,
	swelabAlfaHandler *swelabalfa.Handler,
) *DeviceStrategy {
	return &DeviceStrategy{
		defaultHandler:     defaultHandler,
		analyxTriaHandler:  analyxTriaHandler,
		analyxPancaHandler: analyxPancaHandler,
		swelabAlfaHandler:  swelabAlfaHandler,
	}
}

func (d *DeviceStrategy) ChooseDeviceHandler(device entity.Device) (server.TCPHandler, error) {
	// TODO: Add more device types here and change the default handler to the correct one
	switch device.Type {
	case entity.DeviceTypeBA200, entity.DeviceTypeBA400:
		return d.defaultHandler, nil
	case entity.DeviceTypeA15:
		return nil, entity.ErrDeviceTypeNotSupport
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
