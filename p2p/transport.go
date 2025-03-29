package p2p

import "net"

// Peer is an interface that represents a remote node
type Peer interface {
  net.Conn
  Send([]byte) error
}

// Transport handles communications between nodes over Network
// It can be TCP, UDP, or any other protocol
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
  Dial(string) error
}
