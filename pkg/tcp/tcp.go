package tcp

import (
	"fmt"
	"net"
	"time"
)

// TCP is a struct handle TCP connection
type TCP struct {
	net.Conn
}

// NewTCP returns a new TCP connection
func NewTCP(host string, port int, timeout time.Duration) (*TCP, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d",
		host,
		port,
	), timeout)
	if err != nil {
		return nil, err
	}

	return &TCP{
		Conn: conn,
	}, nil
}

// Send sends TCP message
func (t *TCP) Send(message string) (string, error) {
	_, err := t.Write([]byte(message))
	if err != nil {
		return "", err
	}

	// Read the response from the server
	buf := make([]byte, 4096) // Adjust buffer size as needed
	n, err := t.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Return the response as a string
	return string(buf[:n]), nil
}
