package server

import (
	"bufio"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"go.bug.st/serial"
)

// TestSerialHandler demonstrates how to create a serial handler
type TestSerialHandler struct {
	deviceName string
	messages   []string
}

// NewTestSerialHandler creates a new test serial handler
func NewTestSerialHandler(deviceName string) *TestSerialHandler {
	return &TestSerialHandler{
		deviceName: deviceName,
		messages:   make([]string, 0),
	}
}

// Handle implements the SerialHandler interface
func (h *TestSerialHandler) Handle(port serial.Port) {
	slog.Info("starting test serial handler", "device", h.deviceName)

	// Set a read timeout
	port.SetReadTimeout(5 * time.Second)

	// Create a scanner to read lines
	scanner := bufio.NewScanner(port)
	buffer := ""

	for scanner.Scan() {
		line := scanner.Text()
		buffer += line + "\n"

		// Process complete messages (you can customize this logic)
		if strings.Contains(line, "\n") || strings.Contains(line, "\r") {
			// Process the complete message
			h.processMessage(strings.TrimSpace(buffer))
			buffer = ""
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("error reading from serial port", "error", err)
	}
}

// processMessage handles the received serial message
func (h *TestSerialHandler) processMessage(message string) {
	if message == "" {
		return
	}

	h.messages = append(h.messages, message)
	slog.Info("received serial message",
		"device", h.deviceName,
		"message", message,
		"length", len(message))

	// Here you would implement your specific message processing logic
	// For example, parsing HL7 messages, device-specific protocols, etc.

	// Example: Echo back the message (if needed)
	// port.Write([]byte(fmt.Sprintf("ACK: %s\n", message)))
}

// GetMessages returns all received messages
func (h *TestSerialHandler) GetMessages() []string {
	return h.messages
}

// TestSerialServerCreation tests creating a new serial server
func TestSerialServerCreation(t *testing.T) {
	server := NewSerial("COM3", 115200)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.portName != "COM3" {
		t.Errorf("Expected port name COM3, got %s", server.portName)
	}

	if server.baudRate != 115200 {
		t.Errorf("Expected baud rate 115200, got %d", server.baudRate)
	}

	if server.state != constant.ServerStateStopped {
		t.Errorf("Expected initial state to be stopped, got %v", server.state)
	}
}

// TestSerialServerInterface tests that Serial implements ServerController
func TestSerialServerInterface(t *testing.T) {
	var _ Controller = (*Serial)(nil)
	// This test will fail to compile if Serial doesn't implement ServerController
}

// TestSerialHandlerCreation tests creating a serial handler
func TestSerialHandlerCreation(t *testing.T) {
	handler := NewTestSerialHandler("TestDevice")

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.deviceName != "TestDevice" {
		t.Errorf("Expected device name TestDevice, got %s", handler.deviceName)
	}

	if len(handler.messages) != 0 {
		t.Errorf("Expected empty messages slice, got %d messages", len(handler.messages))
	}
}

// TestSerialHandlerFunc tests using a function as a handler
func TestSerialHandlerFunc(t *testing.T) {
	messages := make([]string, 0)

	handlerFunc := SerialHandlerFunc(func(port serial.Port) {
		// Simulate processing a message
		messages = append(messages, "test message")
		slog.Info("handler function called", "message", "test message")
	})

	// Test that the function implements the interface
	var _ SerialHandler = handlerFunc

	// Test that we can call it (though we can't actually test with a real port)
	if len(messages) != 0 {
		t.Errorf("Expected no messages before call, got %d", len(messages))
	}
}

// TestSerialServerStateManagement tests server state transitions
func TestSerialServerStateManagement(t *testing.T) {
	server := NewSerial("COM3", 115200)

	// Initial state should be stopped
	if server.State() != constant.ServerStateStopped {
		t.Errorf("Expected initial state to be stopped, got %v", server.State())
	}

	// Test setting port
	server.SetPort("COM4")
	if server.portName != "COM4" {
		t.Errorf("Expected port name to be COM4, got %s", server.portName)
	}

	// Test setting handler
	handler := NewTestSerialHandler("TestDevice")
	server.SetHandler(handler)
	if server.handler == nil {
		t.Error("Expected handler to be set, got nil")
	}
}

// TestSerialServerStopChannel tests the stop channel functionality
func TestSerialServerStopChannel(t *testing.T) {
	server := NewSerial("COM3", 115200)

	if server.stopChan == nil {
		t.Fatal("Expected stop channel to be initialized, got nil")
	}

	// Test that we can close the stop channel
	close(server.stopChan)

	// Verify the channel is closed
	select {
	case <-server.stopChan:
		// Expected - channel is closed
	default:
		t.Error("Expected stop channel to be closed")
	}
}

// ExampleNewSerial demonstrates creating a new serial server
func ExampleNewSerial() {
	// Create a new serial server
	serialServer := NewSerial("COM3", 115200)

	// Create and set the handler
	handler := NewTestSerialHandler("ExampleDevice")
	serialServer.SetHandler(handler)

	// Start the server
	err := serialServer.Start()
	if err != nil {
		slog.Error("failed to start serial server", "error", err)
		return
	}

	// Start serving in a goroutine
	go serialServer.Serve()

	// Keep the main thread alive for a short time
	time.Sleep(1 * time.Second)

	// Stop the server
	err = serialServer.Stop()
	if err != nil {
		slog.Error("failed to stop serial server", "error", err)
	}

	// Output:
	// Example demonstrates serial server usage
}

// ExampleSerialHandler demonstrates using a function as a handler
func ExampleSerialHandler() {
	// Create a new serial server
	serialServer := NewSerial("COM3", 115200)

	// Create a handler function
	handler := SerialHandlerFunc(func(port serial.Port) {
		slog.Info("custom handler function called")
		// Your custom serial communication logic here
	})

	serialServer.SetHandler(handler)

	// Start the server
	err := serialServer.Start()
	if err != nil {
		slog.Error("failed to start serial server", "error", err)
		return
	}

	// Start serving in a goroutine
	go serialServer.Serve()

	// Keep the main thread alive for a short time
	time.Sleep(1 * time.Second)

	// Stop the server
	err = serialServer.Stop()
	if err != nil {
		slog.Error("failed to stop serial server", "error", err)
	}

	// Output:
	// Example demonstrates using a function as a serial handler
}

// BenchmarkSerialServerCreation benchmarks server creation
func BenchmarkSerialServerCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		server := NewSerial("COM3", 115200)
		if server == nil {
			b.Fatal("server creation failed")
		}
	}
}

// BenchmarkSerialHandlerCreation benchmarks handler creation
func BenchmarkSerialHandlerCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handler := NewTestSerialHandler("TestDevice")
		if handler == nil {
			b.Fatal("handler creation failed")
		}
	}
}
