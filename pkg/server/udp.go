package server

import (
	"github.com/sirupsen/logrus"
	"net"

	"github.com/anjmao/friends/pkg/types"
)

const udpBufferSize = 65507

func NewUDPServer() Friends {
	return &UDPServer{}
}

type UDPServer struct {
	handler ConnHandler
}

// ListenAndServe starts listening and accepting new UDP packets.
func (s *UDPServer) ListenAndServe(addr string) error {
	if s.handler == nil {
		return errHandlerNotRegistered
	}

	p, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	for {
		buffer := make([]byte, udpBufferSize)
		n, caddr, err := p.ReadFrom(buffer)
		if err != nil {
			logrus.Errorf("could not read packets: %v", err)
			break
		}

		s.handlePacket(p, n, buffer, caddr)
	}
	return nil
}

// Handle registers global handler.
func (s *UDPServer) Handle(handler ConnHandler) {
	s.handler = handler
}

func (s *UDPServer) handlePacket(p net.PacketConn, n int, b []byte, caddr net.Addr) {
	msg := types.DecodeMsg(b[:n])
	s.handler(&ConnContext{udpConn: p, addr:caddr}, msg)
}
