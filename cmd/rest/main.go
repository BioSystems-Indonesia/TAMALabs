package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/app"
	"github.com/oibacidem/lims-hl-seven/internal/util"
)

func main() {
	go main_tcp()
	server := app.InitRestApp()
	go openb()
	server.Serve()
}

func main_tcp() {
	server := app.InitTCPApp()
	server.Serve()
}

func openb() {
	if util.IsDevelopment() {
		return
	}

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
