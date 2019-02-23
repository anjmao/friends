package server

import (
	"errors"

	"github.com/anjmao/friends/pkg/types"
)

// Friends interface describes common server abstraction for TCP/UDP.
type Friends interface {
	ListenAndServe(addr string) error
	Handle(handler ConnHandler)
}

var (
	errHandlerNotRegistered = errors.New("handler is not registered")
)

// ConnHandler abstracts incoming data handling for TCP/UDP protocols.
type ConnHandler func(ctx *ConnContext, msg *types.Msg)
