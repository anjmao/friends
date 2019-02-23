package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/anjmao/friends/pkg/client"
	"github.com/sirupsen/logrus"
)

var (
	protocol = flag.String("protocol", "tcp", "Friends network protocol")
	addr     = flag.String("addr", ":8080", "Server address")
	user     = flag.String("user", "", "User payload")
)

func main() {
	flag.Parse()

	var c client.Friends
	switch *protocol {
	case "tcp":
		c = client.NewTCPClient()
	case "udp":
		c = client.NewUDPClient()
	default:
		logrus.Fatalf("unknown protocol %s", *protocol)
	}

	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		<-gracefulStop
		if err := c.Close(); err != nil {
			logrus.Fatalf("could not close client: %v", err)
		}
		os.Exit(0)
	}()

	if err := c.Connect(*addr, *user); err != nil {
		logrus.Fatalf("could not connect to %s on protocol %s: %v", *addr, *protocol, err)
	}

	go c.PingLoop()
	c.ListenIncoming()
}
