package mllp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
)

const (
	mllpStartBlock     = 0x0b
	mllpEndBlock       = 0x1c
	mllpCarriageReturn = 0x0d
)

// MLLPClient reads and writes MLLP.
type Client struct {
	r *bufio.Reader
	w io.Writer
}

// Dial will connect to the host and return a new MLLPClient.
func Dial(host string) (*Client, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}

	return NewClient(conn), nil
}

// NewMLLPClient returns a new MLLPClient.
func NewClient(c net.Conn) *Client {
	return &Client{bufio.NewReader(c), c}
}

// Write writes a message to MLLP client.
// Returns an error if writing a message fails.
func (c *Client) Write(message []byte) error {
	if _, err := c.w.Write([]byte{mllpStartBlock}); err != nil {
		return err
	}
	if _, err := c.w.Write(message); err != nil {
		return err
	}
	if _, err := c.w.Write([]byte{mllpEndBlock, mllpCarriageReturn}); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReadAll() ([]byte, error) {
	b, err := c.r.ReadByte()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot read the first byte of mllp message, %w", err)
	}
	if b != mllpStartBlock {
		slog.Error("Invalid mllp start block")
		return nil, nil
	}
	payload, err := c.r.ReadBytes(mllpEndBlock)
	if err != nil {
		return nil, fmt.Errorf("cannot read mllp message %w", err)
	}
	// Remove the mllpEndBlock
	payload = payload[:len(payload)-1]
	b, err = c.r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("cannot read the last byte of mllp message, %w", err)
	}
	if b != mllpCarriageReturn {
		return nil, errors.New("mllp: protocol error, missing End Carriage Return")
	}
	return payload, nil
}

func (c *Client) ReadMultiMessage() ([][]byte, error) {
	slog.Info("read multi message")

	var messages [][]byte
	for {
		message, err := c.ReadAll()
		if err != nil {
			slog.Error("error reading mllp message", "error", err)
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}
		if len(message) == 0 {
			break
		}

		slog.Info("read mllp message", "message", string(message))
		messages = append(messages, message)
	}

	return messages, nil
}

func (c *Client) ReadAllRaw() ([]byte, error) {
	var payload []byte
	for {
		b, err := c.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
		payload = append(payload, b)
	}

	return payload, nil
}

// Read implements io.Reader.
func (c *Client) Read(p []byte) (n int, err error) {
	var b []byte

	b, err = c.ReadAll()
	if err != nil {
		return 0, err
	}

	if len(b) > len(p) {
		return 0, errors.New("mllp: message too long")
	}

	n = copy(p, b)
	return n, nil
}
