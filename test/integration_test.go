package test

import (
	"github.com/anjmao/friends/pkg/client"
	"github.com/anjmao/friends/pkg/server"
	"testing"
	"time"
)

const (
	serveAddr          = ":9090"
	testWaitTime       = 10 * time.Millisecond
	checkStateInterval = 200 * time.Millisecond
)

type ClientFunc = func() client.Friends
type ServerFunc = func() server.Friends

func TestServerAndClientConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tests := []struct {
		name       string
		serverFunc ServerFunc
		clientFunc ClientFunc
	}{
		{
			name:"TCP server with TCP client",
			serverFunc: func() server.Friends { return server.NewTCPServer() },
			clientFunc: func() client.Friends { return client.NewTCPClient() },
		},
		{
			name:"UDP server with UDP client",
			serverFunc: func() server.Friends { return server.NewUDPServer() },
			clientFunc: func() client.Friends { return client.NewUDPClient() },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			testServerClient(tt, test.clientFunc, test.serverFunc)
		})
	}
}

func testServerClient(t *testing.T, clientFunc ClientFunc, serverFunc ServerFunc) {
	hub := server.NewHub()
	checkTicker := time.NewTicker(checkStateInterval)
	defer checkTicker.Stop()
	done := make(chan struct{})
	go hub.Run(checkTicker.C, done)

	var clients []client.Friends

	users := [...]string{
		"{\"user_id\":1, \"friends\": [2, 3, 4]}",
		"{\"user_id\":2, \"friends\": [1]}",
		"{\"user_id\":3, \"friends\": [2, 3]}",
		"{\"user_id\":4, \"friends\": [3]}",
	}

	// Start server.
	srv := serverFunc()
	srv.Handle(hub.IncomingMessageHandler)
	go func() {
		if err := srv.ListenAndServe(serveAddr); err != nil {
			t.Fatal(err)
		}
	}()

	// Wait for server to start.
	time.Sleep(testWaitTime)

	// Cleanup clients.
	defer func() {
		for _, c := range clients {
			if err := c.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}()

	// Connect clients.
	for _, u := range users {
		c := clientFunc()
		if err := c.Connect(serveAddr, u); err != nil {
			t.Fatal(err)
		}
		clients = append(clients, c)

		go c.PingLoop()
		go c.ListenIncoming()

		// Wait some time for clients to finish connecting.
		time.Sleep(10 * time.Millisecond)
	}

	done <- struct{}{}

	expectedUsersLen := len(users)
	actualUsersLen := len(hub.Users())
	if expectedUsersLen != actualUsersLen {
		t.Fatalf("expected to have %d users, got %d", expectedUsersLen, actualUsersLen)
	}
}
