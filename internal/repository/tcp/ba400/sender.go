package ba400

import (
	"fmt"
	"net"

	"github.com/oibacidem/lims-hl-seven/pkg/mllp"
)

type Sender struct {
	host string
}

// SendAlpha1 is one of interface that I want to use to send
func (s *Sender) SendRaw(msg []byte) ([]byte, error) {
	conn, err := net.Dial("tcp", s.host)
	if err != nil {
		return nil, fmt.Errorf("connect error: %w", err)
	}
	defer conn.Close()

	c := mllp.NewClient(conn)
	err = c.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("send/write error: %w", err)
	}

	res, err := c.ReadAll()
	if err != nil {
		return res, fmt.Errorf("recv/read error: %w", err)
	}

	return res, nil
}

