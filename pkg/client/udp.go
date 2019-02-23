package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/anjmao/friends/pkg/types"

	"github.com/sirupsen/logrus"
)

const (
	updPingInterval = 100 * time.Millisecond
	udpBufferSize   = 65507
)

func NewUDPClient() Friends {
	return &UDPClient{}
}

// UDPClient implements Friends using UDP protocol.
type UDPClient struct {
	conn   *net.UDPConn
	userID int
}

func (c *UDPClient) Connect(addr, user string) error {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	c.conn = conn

	req := &types.LoginRequest{}
	r := strings.NewReader(user)
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return fmt.Errorf("invalid login payload: %v", err)
	}

	if err := c.sendMessage(types.CmdLogin, req); err != nil {
		return fmt.Errorf("could send data: %v", err)
	}

	c.userID = req.UserID
	// While for TCP client we know if request was sent successfully
	// with UDP we could not be sure that packet was not lost.
	// Ideally UDP client should listen to some ACK message before
	// starting the game, but here we keep it simple and assume that
	// client sent initial request and server got it.
	fmt.Printf("user %d connected to the game\n", c.userID)
	return nil
}

func (c *UDPClient) ListenIncoming() {
	for {
		buffer := make([]byte, udpBufferSize)
		n, _, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			logrus.Errorf("could not read packet: %v", err)
			return
		}

		msg := types.DecodeMsg(buffer[:n])

		switch msg.Cmd {
		case types.CmdStatusChange:
			logrus.Infof("friend status changed: %s", msg.Data)
		}
	}
}

func (c *UDPClient) Close() error {
	if c.conn == nil {
		return c.conn.Close()
	}
	return nil
}

func (c *UDPClient) PingLoop() {
	for {
		time.Sleep(updPingInterval)
		err := c.sendMessage(types.CmdPing, &types.PingRequest{UserID: c.userID})
		if err != nil {
			logrus.Errorf("could not ping server: %v", err)
		}
	}
}

func (c *UDPClient) sendMessage(cmd types.CommandType, v interface{}) error {
	msg, err := types.EncodeMsg(cmd, v)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(msg)
	return err
}
