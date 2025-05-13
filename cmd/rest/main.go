package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/app"
	"github.com/oibacidem/lims-hl-seven/internal/util"
)

func main() {
	//l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ AddSource: true, }))
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(l)

	server := app.InitRestApp()
	go openb()
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
		log.Println(err)
	}
}
