package client

// Friends interface describe common client abstraction over TCP/UDP.
type Friends interface {
	Connect(addr, user string) error
	PingLoop()
	ListenIncoming()
	Close() error
}
