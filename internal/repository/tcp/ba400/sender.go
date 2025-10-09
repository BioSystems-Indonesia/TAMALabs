package ba400

import (
	"fmt"
	"net"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/pkg/mllp"
)

type Sender struct {
	host     string
	deadline time.Duration
}

// SendAlpha1 is one of interface that I want to use to send
func (s *Sender) SendRaw(msg []byte) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", s.host, s.deadline)
	if err != nil {
		return nil, fmt.Errorf("connect error: %w", err)
	}
	defer conn.Close()

	if s.deadline > 0 {
		err := conn.SetDeadline(time.Now().Add(s.deadline))
		if err != nil {
			return nil, fmt.Errorf("set deadline error: %w", err)
		}
	}

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
