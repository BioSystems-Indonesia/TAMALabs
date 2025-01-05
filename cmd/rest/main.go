package main

import (
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/app"
)

func main() {
	go main_tcp()
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	server := app.InitRestApp(&cfg)
	server.Serve()
}

func main_tcp() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	server := app.InitTCPApp(&cfg)
	server.Serve()
}
