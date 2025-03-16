package p2p

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	shakeHands    HandshakeFunc
	mu            sync.RWMutex
	peers         map[net.Addr]Peer
  decoder       Decoder
}

type Temp struct {}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound, // outbound OR inbound
	}
}

func NewTCPTransport(listenAddress string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddress,
		shakeHands:    NOPHandshakeFunc,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}
	go t.StartAcceptLoop()
	return nil
}

func (t *TCPTransport) StartAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.HandleConn(conn)
	}
}

func (t *TCPTransport) HandleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	fmt.Printf("New Incoming connection %v\n", peer)

	if err := t.shakeHands(peer); err != nil {
		fmt.Printf("Handshake failed: %s\n", err)
		return
	}

  msg := &Temp{}
  for {
    if err := t.decoder.Decode(conn, msg); err != nil {
      fmt.Printf("Decoding Error: %s\n", err)
    }
    
  }
}
