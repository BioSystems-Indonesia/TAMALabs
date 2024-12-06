package main

import (
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/app"
)

func main() {
	server := app.InitRestApp(&config.Schema{})
	server.Serve()
}
