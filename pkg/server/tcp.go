package server

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
)

// TCPHandler is an interface for TCP handler
type TCPHandler interface {
	Handle(conn *net.TCPConn)
}

type TCPHandlerFunc func(c *net.TCPConn)

func (t TCPHandlerFunc) Handle(c *net.TCPConn) {
	t(c)
}

// TCP structure
type TCP struct {
	port            string
	listener        *net.TCPListener
	handler         TCPHandler
	state           constant.ServerState
	timeLastConnect time.Time
}

// Ensure TCP implements ServerController interface
var _ Controller = (*TCP)(nil)

// NewTCP returns a new TCP server
func NewTCP(port string) *TCP {
	return &TCP{
		port: port,
	}
}

func (t *TCP) SetPort(port string) {
	t.port = port
}

func (t *TCP) SetHandler(handler TCPHandler) {
	t.handler = handler
}

// Accept returns the TCP client
func (t *TCP) Accept() *net.TCPConn {
	conn, err := t.listener.AcceptTCP()
	if err != nil {
		slog.Error("error accepting connection", "error", err)
	}
	return conn
}

// Start
func (t *TCP) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", ":"+t.port)
	if err != nil {
		return fmt.Errorf("error resolving address :%s : %w", t.port, err)
	}

	t.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("error listen tcp: %w", err)
	}

	slog.Info("start tcp server", "addr", addr)

	return nil
}

// Serve
func (t *TCP) Serve() {
	t.state = constant.ServerStateServing
	defer func() { t.state = constant.ServerStateStopped }()
	for {
		if t.state == constant.ServerStateStopped {
			slog.Info("server stopped", "state", t.state, "port", t.port)
			return
		}

		if t.state == constant.ServerStateConnect {
			if time.Since(t.timeLastConnect) > disconnectTimeout {
				t.state = constant.ServerStateServing
				slog.Info("disconnect timeout, change to serving")
			}
		}

		conn, err := t.listener.AcceptTCP()
		if err != nil {
			slog.Error("error accepting connection", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		t.timeLastConnect = time.Now()
		t.state = constant.ServerStateConnect

		t.handler.Handle(conn)
	}
}

// State
func (t *TCP) State() constant.ServerState {
	return t.state
}

// Stop
func (t *TCP) Stop() error {
	// TODO try to do wait of unfinished handle func
	t.state = constant.ServerStateStopped
	return t.listener.Close()
}
