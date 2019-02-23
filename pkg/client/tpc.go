package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/anjmao/friends/pkg/types"
	"github.com/sirupsen/logrus"
)

const (
	tcpPingInterval = 100 * time.Millisecond
)

func NewTCPClient() Friends {
	return &TCPClient{}
}

// TCPClient implements Friends using TCP protocol.
type TCPClient struct {
	conn   net.Conn
	userID int
}

func (c *TCPClient) Connect(addr, user string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("could open %s connection on addr %s: %v", "tcp", addr, err)
	}
	c.conn = conn

	req := &types.LoginRequest{}
	r := strings.NewReader(user)
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return fmt.Errorf("invalid user login payload: %v", err)
	}

	if err := c.sendMessage(types.CmdLogin, req); err != nil {
		return fmt.Errorf("could send data: %v", err)
	}
	c.userID = req.UserID
	logrus.Infof("user %d connected to the game\n", c.userID)
	return nil
}

func (c *TCPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *TCPClient) ListenIncoming() {
	scanner := bufio.NewScanner(c.conn)
	for {
		if ok := scanner.Scan(); !ok {
			return
		}

		msg := types.DecodeMsg(scanner.Bytes())
		switch msg.Cmd {
		case types.CmdStatusChange:
			logrus.Infof("friend status changed: %s", msg.Data)
		}
	}
}

func (c *TCPClient) PingLoop() {
	for {
		time.Sleep(tcpPingInterval)
		err := c.sendMessage(types.CmdPing, &types.PingRequest{UserID: c.userID})
		if err != nil {
			logrus.Errorf("could not ping server: %v", err)
		}
	}
}

func (c *TCPClient) sendMessage(cmd types.CommandType, v interface{}) error {
	msg, err := types.EncodeMsg(cmd, v)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(msg)
	return err
}
