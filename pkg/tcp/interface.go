package tcp

// TCPer is a interface handle TCP connection
type TCPer interface {
	Send(message string) (string, error)
}
