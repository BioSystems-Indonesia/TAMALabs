package server

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
)

// TCPHandler is an interface for TCP handler
type TCPHandler interface {
	Handle(conn *net.TCPConn)
}

type TCPHandlerFunc func(c *net.TCPConn)

func (t TCPHandlerFunc) Handle(c *net.TCPConn) {
	t(c)
}

// ConnectionInfo tracks individual connection information
type ConnectionInfo struct {
	conn         *net.TCPConn
	lastActivity time.Time
	isActive     bool
}

// TCP structure
type TCP struct {
	port        string
	listener    *net.TCPListener
	handler     TCPHandler
	state       constant.ServerState
	connections map[string]*ConnectionInfo
	mu          sync.RWMutex
}

// Ensure TCP implements ServerController interface
var _ Controller = (*TCP)(nil)

// NewTCP returns a new TCP server
func NewTCP(port string) *TCP {
	return &TCP{
		port:        port,
		connections: make(map[string]*ConnectionInfo),
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

// cleanupDeadConnections removes connections that are no longer active
func (t *TCP) cleanupDeadConnections() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for addr, connInfo := range t.connections {
		if !connInfo.isActive || time.Since(connInfo.lastActivity) > disconnectTimeout {
			slog.Info("removing dead connection", "address", addr)
			if connInfo.conn != nil {
				connInfo.conn.Close()
			}
			delete(t.connections, addr)
		}
	}

	// Update server state based on active connections
	if len(t.connections) > 0 {
		t.state = constant.ServerStateConnect
	} else {
		t.state = constant.ServerStateServing
	}
}

// Serve
func (t *TCP) Serve() {
	t.state = constant.ServerStateServing
	defer func() { t.state = constant.ServerStateStopped }()

	// Start a goroutine to periodically cleanup dead connections
	cleanupTicker := time.NewTicker(2 * time.Second)
	defer cleanupTicker.Stop()

	go func() {
		for range cleanupTicker.C {
			if t.state == constant.ServerStateStopped {
				return
			}
			t.cleanupDeadConnections()
		}
	}()

	for {
		if t.state == constant.ServerStateStopped {
			slog.Info("server stopped", "state", t.state, "port", t.port)
			return
		}

		conn, err := t.listener.AcceptTCP()
		if err != nil {
			slog.Error("error accepting connection", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Add new connection to tracking
		remoteAddr := conn.RemoteAddr().String()
		t.mu.Lock()
		t.connections[remoteAddr] = &ConnectionInfo{
			conn:         conn,
			lastActivity: time.Now(),
			isActive:     true,
		}
		t.state = constant.ServerStateConnect
		t.mu.Unlock()

		slog.Info("new connection accepted", "address", remoteAddr)

		// Handle connection in a separate goroutine
		go t.handleConnection(conn, remoteAddr)
	}
}

// handleConnection handles individual connection
func (t *TCP) handleConnection(conn *net.TCPConn, remoteAddr string) {
	defer func() {
		conn.Close()
		t.mu.Lock()
		delete(t.connections, remoteAddr)
		t.mu.Unlock()
		slog.Info("connection closed", "address", remoteAddr)
	}()

	// Update last activity
	t.mu.Lock()
	if connInfo, exists := t.connections[remoteAddr]; exists {
		connInfo.lastActivity = time.Now()
	}
	t.mu.Unlock()

	// Let the handler handle the connection
	t.handler.Handle(conn)

	// Mark connection as inactive when handler returns
	t.mu.Lock()
	if connInfo, exists := t.connections[remoteAddr]; exists {
		connInfo.isActive = false
	}
	t.mu.Unlock()
}

// State
func (t *TCP) State() constant.ServerState {
	return t.state
}

// Stop
func (t *TCP) Stop() error {
	t.state = constant.ServerStateStopped

	// Close all active connections
	t.mu.Lock()
	for addr, connInfo := range t.connections {
		if connInfo.conn != nil {
			slog.Info("closing connection during stop", "address", addr)
			connInfo.conn.Close()
		}
	}
	// Clear connections map
	t.connections = make(map[string]*ConnectionInfo)
	t.mu.Unlock()

	// Close listener
	if t.listener != nil {
		return t.listener.Close()
	}
	return nil
}
