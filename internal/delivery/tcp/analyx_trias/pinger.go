package analyxtrias

import (
	"bufio"
	"log/slog"
	"net"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
)

type Pinger struct {
	conn *net.TCPConn
}

func NewPinger(conn *net.TCPConn) *Pinger {
	return &Pinger{conn: conn}
}

func (p *Pinger) IsPing() bool {
	bufReader := bufio.NewReader(p.conn)
	b, err := bufReader.Peek(2)
	if err != nil {
		slog.Error("error on ping peek", "err", err.Error())
		return false
	}

	if len(b) == 0 {
		return false
	}

	if b[0] == constant.StartOfText && b[1] == constant.StartOfText {
		return true
	}

	return false
}

func (p *Pinger) ReturnACK() (int, error) {
	return p.conn.Write([]byte{constant.Acknowledge})
}
