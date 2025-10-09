package server

import (
	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/server"
)

type DummyServer struct {
}

// DummyServer is a dummy server that is used to represent a server that is not running
// because the device does not need the tcp handler
func NewDummyServer() *DummyServer {
	return &DummyServer{}
}

var _ server.Controller = &DummyServer{}

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
