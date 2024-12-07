package server

import (
	"log"
	"net"
	"os"
	"os/signal"
)

type TCPServer interface {
	Serve()
	GetClient() *net.TCPConn
}

// TCP structure
type TCP struct {
	port     string
	listener *net.TCPListener
}

// NewTCP returns a new TCP server
func NewTCP(port string) TCPServer {
	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error resolving address: %v", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Printf("Listener created on %s", listener.Addr().String())
	return &TCP{
		port:     port,
		listener: listener,
	}
}

// GetClient returns the TCP client
func (t *TCP) GetClient() *net.TCPConn {
	conn, err := t.listener.AcceptTCP()
	if err != nil {
		log.Fatalf("Error accepting connection: %v", err)
	}
	return conn
}

// Serve starts serving incoming connections
func (t *TCP) Serve() {
	log.Printf("Server started on port %s", t.port)

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	if err := t.listener.Close(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
