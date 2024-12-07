package main

import (
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/app"
)

func main() {
	server := app.InitTCPApp(&config.Schema{})
	server.Serve()
}
