package main

import (
	"github.com/oibacidem/lims-hl-seven/internal/app"
)

func main() {
	server := app.InitTCPApp()
	server.Serve()
}
