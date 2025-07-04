package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"os/exec"
	"runtime"
	"time"

	"github.com/energye/systray"
	"github.com/oibacidem/lims-hl-seven/internal/app"
	"github.com/oibacidem/lims-hl-seven/internal/util"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

var Version = ""

//go:embed trayicon.ico
var trayicon []byte

func main() {
	//l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ AddSource: true, }))
	//l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	l := slog.Default()

	slog.SetDefault(l)

	log.Println("version: ", Version)

	server := app.InitRestApp()

	go openb()
	go opensystray(server)
	server.Serve()
}

func openb() {
	if Version == "" || util.IsDevelopment() {
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
		log.Println(err)
	}
}

func opensystray(server server.RestServer) {
	systray.Run(func() {
		systray.SetTitle("LIMS HL Seven")
		systray.SetTooltip("LIMS HL Seven")
		systray.SetIcon(trayicon)

		systray.AddMenuItem("Open Browser", "Open Browser").Click(func() {
			openbrowser("http://127.0.0.1:8322")
		})

		systray.AddMenuItem("Quuit", "Stop Server").Click(func() {
			if err := server.Stop(); err != nil {
				log.Println("Error stopping server:", err)
			} else {
				log.Println("Server stopped successfully")
			}
			systray.Quit()
		})
	}, func() {})
}
