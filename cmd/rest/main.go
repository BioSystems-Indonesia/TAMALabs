package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

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
	go openb()
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

func openb() {
	time.Sleep(3 * time.Second)
	openbrowser("http://127.0.0.1:8322")
}


func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
