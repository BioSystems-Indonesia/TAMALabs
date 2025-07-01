package server

import (
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
)

type DummyServer struct {
}

// DummyServer is a dummy server that is used to represent a server that is not running
// because the device does not need the tcp handler
func NewDummyServer() *DummyServer {
	return &DummyServer{}
}

var _ repository.TCPServerController = &DummyServer{}

func (d *DummyServer) SetPort(_ string) {
}

func (d *DummyServer) Start() error {
	return nil
}

func (d *DummyServer) State() constant.ServerState {
	return constant.ServerStateNoServer
}

func (d *DummyServer) Serve() {

}

func (d *DummyServer) Stop() error {
	return nil
}
