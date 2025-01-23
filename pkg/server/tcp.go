package server

import (
	"fmt"
	"log/slog"
	"net"
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
	port     string
	listener *net.TCPListener
	handler  TCPHandler
	serving  bool
}

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
	t.serving = true
	defer func() { t.serving = false }()
	for {
		conn, err := t.listener.AcceptTCP()
		if err != nil {
			slog.Error("error accepting connection", "error", err)
			continue
		}
		t.handler.Handle(conn)
		conn.Close()
	}
}

// State
func (t *TCP) State() string {
	if t.serving {
		return "serving"
	}
	return "stoped"
}

// Stop
func (t *TCP) Stop() error {
	// TODO try to do wait of unfinished handle func
	return t.listener.Close()
}
