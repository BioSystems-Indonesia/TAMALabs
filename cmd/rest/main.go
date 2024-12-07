package main

import (
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	server := app.InitRestApp(&cfg)
	server.Serve()
}
