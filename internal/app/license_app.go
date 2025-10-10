package app

import (
	licenserepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/license"
	licenseuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/license"
	"github.com/google/wire"
)

// licenseSet provides constructors for the license service and its filesystem adapters.
var licenseSet = wire.NewSet(
	licenserepo.NewFSKeyLoader,
	licenserepo.NewFSFileLoader,
	provideLicenseService,
	wire.Bind(new(licenseuc.PublicKeyLoader), new(*licenserepo.FSKeyLoader)),
	wire.Bind(new(licenseuc.LicenseFileLoader), new(*licenserepo.FSFileLoader)),
)
