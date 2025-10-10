// Package app is a package that handle dependency injection using wire.

//go:build wireinject
// +build wireinject

package app

import (
	khanzauc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/khanza"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/server"
	"github.com/google/wire"
)

// InitRestApp is a Wire provider function that returns a RestServer.
func InitRestApp() server.RestServer {
	wire.Build(restAppSet)
	return &server.Rest{}
}

// InitCanalHandler is a Wire provider function that returns a CanalHandler.
func InitCanalHandler() *khanzauc.CanalHandler {
	wire.Build(canalHandlerSet)
	return &khanzauc.CanalHandler{}
}
