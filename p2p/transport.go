package p2p

// Peer represents a node in the network
type Peer interface {
  Close() error
}

// Transport handles communications between nodes over Network
// It can be TCP, UDP, or any other protocol
type Transport interface {
	ListenAndAccept() error
  Consume() <-chan RPC
}
