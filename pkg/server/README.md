# Server Package

This package provides server abstractions for both TCP and Serial communication, implementing a common `ServerController` interface for easy interchangeability.

## Overview

The server package includes:

- `TCP` server for TCP/IP communication
- `Serial` server for serial port communication
- Common `ServerController` interface
- Example implementations and usage patterns

## ServerController Interface

Both TCP and Serial servers implement the `ServerController` interface:

```go
type ServerController interface {
    SetPort(port string)
    Start() error
    State() constant.ServerState
    Serve()
    Stop() error
}
```

## Serial Server Usage

### Basic Usage

```go
import "github.com/BioSystems-Indonesia/TAMALabs/pkg/server"

// Create a new serial server
serialServer := server.NewSerial("COM3", 115200)

// Create and set a handler
handler := &MySerialHandler{}
serialServer.SetHandler(handler)

// Start the server
err := serialServer.Start()
if err != nil {
    log.Fatal(err)
}

// Start serving in a goroutine
go serialServer.Serve()

// Later, stop the server
err = serialServer.Stop()
```

### Creating a Custom Handler

Implement the `SerialHandler` interface:

```go
type MySerialHandler struct {
    deviceName string
}

func (h *MySerialHandler) Handle(port serial.Port) {
    // Your serial communication logic here
    // Read from port, process data, send responses, etc.
}
```

### Handler Function

You can also use a function as a handler:

```go
handler := server.SerialHandlerFunc(func(port serial.Port) {
    // Your serial communication logic here
})
serialServer.SetHandler(handler)
```

## TCP Server Usage

The TCP server works similarly:

```go
// Create a new TCP server
tcpServer := server.NewTCP("8080")

// Create and set a handler
handler := &MyTCPHandler{}
tcpServer.SetHandler(handler)

// Start the server
err := tcpServer.Start()
if err != nil {
    log.Fatal(err)
}

// Start serving in a goroutine
go tcpServer.Serve()
```

## Server States

The servers can be in the following states:

- `ServerStateServing`: Server is running and waiting for connections
- `ServerStateConnect`: Server is handling a connection
- `ServerStateStopped`: Server has been stopped
- `ServerStateNoServer`: Server is not running (dummy server)

## Configuration

### Serial Server Configuration

- **Port Name**: The serial port name (e.g., "COM3" on Windows, "/dev/ttyUSB0" on Linux)
- **Baud Rate**: Communication speed (common values: 9600, 19200, 38400, 57600, 115200)
- **Data Bits**: 8 (standard)
- **Parity**: NoParity (standard)
- **Stop Bits**: OneStopBit (standard)

### TCP Server Configuration

- **Port**: The TCP port number as a string

## Error Handling

Both servers provide proper error handling:

- Connection errors are logged
- Graceful shutdown is supported
- State management for monitoring server status

## Example Implementation

See `serial_example.go` for a complete example of how to implement a serial handler and use the serial server.

## Interchangeability

Since both servers implement the same interface, you can easily swap between TCP and Serial servers in your application:

```go
var server server.ServerController

if useSerial {
    server = server.NewSerial("COM3", 115200)
} else {
    server = server.NewTCP("8080")
}

// Use the same interface methods regardless of server type
server.SetHandler(handler)
server.Start()
go server.Serve()
```
