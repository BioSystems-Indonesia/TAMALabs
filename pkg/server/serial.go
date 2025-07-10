package server

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"go.bug.st/serial"
)

// SerialHandler is an interface for Serial handler
type SerialHandler interface {
	Handle(port serial.Port)
}

type SerialHandlerFunc func(port serial.Port)

func (t SerialHandlerFunc) Handle(port serial.Port) {
	t(port)
}

// Serial structure
type Serial struct {
	portName        string
	baudRate        int
	port            serial.Port
	handler         SerialHandler
	state           constant.ServerState
	timeLastConnect time.Time
	stopChan        chan struct{}
}

// Ensure Serial implements ServerController interface
var _ Controller = (*Serial)(nil)

// NewSerial returns a new Serial server
func NewSerial(portName string, baudRate int) *Serial {
	return &Serial{
		portName: portName,
		baudRate: baudRate,
		stopChan: make(chan struct{}),
	}
}

func (t *Serial) SetPort(portName string) {
	t.portName = portName
}

func (t *Serial) SetHandler(handler SerialHandler) {
	t.handler = handler
}

// Start initializes the serial port
func (t *Serial) Start() error {
	mode := &serial.Mode{
		BaudRate: t.baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(t.portName, mode)
	if err != nil {
		return fmt.Errorf("error opening serial port %s: %w", t.portName, err)
	}

	t.port = port
	slog.Info("start serial server", "port", t.portName, "baud_rate", t.baudRate)

	return nil
}

// Serve handles serial communication
func (t *Serial) Serve() {
	t.state = constant.ServerStateServing
	defer func() {
		t.state = constant.ServerStateStopped
		if t.port != nil {
			t.port.Close()
		}
	}()

	for {
		select {
		case <-t.stopChan:
			slog.Info("server stopped", "state", t.state, "port", t.portName)
			return
		default:
			if t.state == constant.ServerStateStopped {
				slog.Info("server stopped", "state", t.state, "port", t.portName)
				return
			}

			if t.state == constant.ServerStateConnect {
				if time.Since(t.timeLastConnect) > disconnectTimeout {
					t.state = constant.ServerStateServing
					slog.Info("disconnect timeout, change to serving")
				}
			}

			// Handle the serial connection
			if t.handler != nil {
				t.timeLastConnect = time.Now()
				t.state = constant.ServerStateConnect
				t.handler.Handle(t.port)
			}

			// Small delay to prevent busy waiting
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// State returns the current server state
func (t *Serial) State() constant.ServerState {
	return t.state
}

// Stop stops the serial server
func (t *Serial) Stop() error {
	t.state = constant.ServerStateStopped
	close(t.stopChan)

	if t.port != nil {
		return t.port.Close()
	}
	return nil
}
