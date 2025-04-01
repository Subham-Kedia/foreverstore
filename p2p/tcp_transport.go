package p2p

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Subham-Kedia/foreverstore/message"
)

type TCPPeer struct {
	net.Conn
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan message.Message

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func (p *TCPPeer) Send(data []byte) error {
	_, err := p.Conn.Write(data)
	return err
}

func (p *TCPPeer) IsOutbound() bool {
	return p.outbound
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound, // outbound OR inbound
	}
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan message.Message),
	}
}

// creating a listener
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	fmt.Printf("Creating listener for %s\n", t.ListenAddr)
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	log.Printf("server is listening on %v\n", t.ListenAddr)
	go t.StartAcceptLoop()
	return nil
}

// accepting incoming connections in a loop
func (t *TCPTransport) StartAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		log.Printf("Incoming Request from %v on %v\n", conn.RemoteAddr(), t.ListenAddr)
		go t.HandleConn(conn, false)
	}
}

func (t *TCPTransport) HandleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("Closing connection: %s\n", err)
		if err != nil {
			conn.Close()
		}
	}()
	fmt.Printf("Handling incoming connection from %s\n", conn.RemoteAddr())
	peer := NewTCPPeer(conn, outbound)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("Handshake failed: %s\n", err)
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			fmt.Printf("OnPeer error: %s\n", err)
			conn.Close()
			return
		}
	}

	// rpc := RPC{}
	var m message.Message
	for {
		err := gob.NewDecoder(conn).Decode(&m)
		if errors.Is(err, net.ErrClosed) {
			fmt.Printf("Connection closed: %s\n", err)
			return
		}
		if err != nil {
			fmt.Printf("Invalid Input: %s\n", err)
			continue
		}
		// m.From = conn.RemoteAddr()
		t.rpcch <- m
	}
}

// Consume implements the Transport Interface
func (t *TCPTransport) Consume() <-chan message.Message {
	return t.rpcch
}

// Close implements the Transport Interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements the Transport Interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	// outbound connection
	go t.HandleConn(conn, true)
	return nil
}
