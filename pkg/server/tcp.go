package server

import (
	"bufio"
	"net"

	"github.com/anjmao/friends/pkg/types"
	"github.com/sirupsen/logrus"
)

func NewTCPServer() Friends {
	return &TCPServer{}
}

type TCPServer struct {
	handler ConnHandler
}

// ListenAndServe starts listening and accepting new TCP connections.
func (s *TCPServer) ListenAndServe(addr string) error {
	if s.handler == nil {
		return errHandlerNotRegistered
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Errorf("failed to accept new conn: %v", err)
			break
		}
		go s.handleConnection(conn)
	}
	return nil
}

// Handle registers global handler.
func (s *TCPServer) Handle(handler ConnHandler) {
	s.handler = handler
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for {
		if ok := scanner.Scan(); !ok {
			return
		}

		msg := types.DecodeMsg(scanner.Bytes())
		s.handler(&ConnContext{tcpConn: conn}, msg)
	}
}
