package p2p

import (
	"net"

	"github.com/Subham-Kedia/foreverstore/message"
)

// Peer is an interface that represents a remote node
type Peer interface {
	net.Conn
	Send([]byte) error
	IsOutbound() bool
}

// Transport handles communications between nodes over Network
// It can be TCP, UDP, or any other protocol
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan message.Message
	Close() error
	Dial(string) error
}
