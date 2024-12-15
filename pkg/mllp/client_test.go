package mllp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Write_Read(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)

		c := NewClient(conn)
		b, err := c.ReadAll()
		require.NoError(t, err)

		err = c.Write(b)
		require.NoError(t, err)
	}()

	c, err := Dial(l.Addr().String())
	require.NoError(t, err)

	msg := []byte("hello, world")

	err = c.Write(msg)
	require.NoError(t, err)

	b := make([]byte, len(msg))
	n, err := c.Read(b)
	require.NoError(t, err)
	assert.Equal(t, len(msg), n)
}
