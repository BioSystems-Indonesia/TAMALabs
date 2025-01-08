// Package app is a package that handle dependency injection using wire.

//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

// InitRestApp is a Wire provider function that returns a RestServer.
func InitRestApp() server.RestServer {
	wire.Build(restAppSet)
	return &server.Rest{}
}

func InitTCPApp() server.TCPServer {
	wire.Build(tcpAppSet)
	return &server.TCP{}
}
