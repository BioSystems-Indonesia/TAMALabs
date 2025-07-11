// Package app is a package that handle dependency injection using wire.

//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	ncc3300 "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/ncc_3300"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

// InitRestApp is a Wire provider function that returns a RestServer.
func InitRestApp() server.RestServer {
	wire.Build(restAppSet)
	return &server.Rest{}
}

// InitSerialNCC3300App is a Wire provider function that returns a serial NCC3300 Handler.
func InitSerialNCC3300App() *ncc3300.Handler {
	wire.Build(serialNCC3300AppSet)
	return &ncc3300.Handler{}
}
