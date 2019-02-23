package main

import (
	"flag"
	"time"

	"github.com/anjmao/friends/pkg/server"
	"github.com/sirupsen/logrus"
)

var (
	protocol = flag.String("protocol", "tcp", "Friends server protocol")
	addr     = flag.String("addr", ":8080", "Serve address")
)

func main() {
	flag.Parse()

	hub := server.NewHub()
	checkTicker := time.NewTicker(server.CheckUsersStateInterval)
	defer checkTicker.Stop()
	done := make(chan struct{})
	go hub.Run(checkTicker.C, done)

	var srv server.Friends
	switch *protocol {
	case "tcp":
		srv = server.NewTCPServer()
	case "udp":
		srv = server.NewUDPServer()
	default:
		logrus.Fatalf("unknown protocol %s", *protocol)
	}

	srv.Handle(hub.IncomingMessageHandler)
	if err := srv.ListenAndServe(*addr); err != nil {
		done<- struct{}{}
		logrus.Fatal(err)
	}
}
