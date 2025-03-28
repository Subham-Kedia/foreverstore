package p2p

import (
	"errors"
	"fmt"
	"net"
	"sync"
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
	rpcch    chan RPC

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func (p *TCPPeer) Send(data []byte) error {
	_, err := p.Conn.Write(data)
	return err
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
		rpcch:            make(chan RPC),
	}
}

// creating a listener
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
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

	peer := NewTCPPeer(conn, outbound)

	fmt.Printf("New Inbound connection %v %v\n", peer.Conn.LocalAddr(), peer.Conn.RemoteAddr())

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

	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if errors.Is(err, net.ErrClosed) {
			fmt.Printf("Connection closed: %s\n", err)
			return
		}
		if err != nil {
			fmt.Printf("Invalid Input: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}

// Consume implements the Transport Interface
func (t *TCPTransport) Consume() <-chan RPC {
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
