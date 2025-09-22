package app

import (
	"github.com/google/wire"
	licenserepo "github.com/oibacidem/lims-hl-seven/internal/repository/license"
	licenseuc "github.com/oibacidem/lims-hl-seven/internal/usecase/license"
)

// licenseSet provides constructors for the license service and its filesystem adapters.
var licenseSet = wire.NewSet(
	licenserepo.NewFSKeyLoader,
	licenserepo.NewFSFileLoader,
	provideLicenseService,
	wire.Bind(new(licenseuc.PublicKeyLoader), new(*licenserepo.FSKeyLoader)),
	wire.Bind(new(licenseuc.LicenseFileLoader), new(*licenserepo.FSFileLoader)),
)
