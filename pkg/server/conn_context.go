package server

import (
	"fmt"
	"net"
	"time"
)

const writeDeadline = 3 * time.Second

// ConnContext is a wrapper around TCP/UDP which helps
// to abstract game logic and make it independent from the
// communication protocol.
type ConnContext struct {
	tcpConn net.Conn
	udpConn net.PacketConn
	addr    net.Addr
}

// write writes data to underlying network connection.
func (c *ConnContext) write(b []byte) error {
	if c.tcpConn != nil {
		if err := c.tcpConn.SetWriteDeadline(time.Now().Add(writeDeadline)); err != nil {
			return fmt.Errorf("could not set write deadline: %v", err)
		}
		_, err := c.tcpConn.Write(b)
		return err
	}

	if err := c.udpConn.SetWriteDeadline(time.Now().Add(writeDeadline)); err != nil {
		return fmt.Errorf("could not set write deadline: %v", err)
	}
	_, err := c.udpConn.WriteTo(b, c.addr)
	return err
}

// close closes TCP connections. Does nothing for UDP.
func (c *ConnContext) close() error {
	if c.tcpConn != nil {
		return c.tcpConn.Close()
	}
	return nil
}
