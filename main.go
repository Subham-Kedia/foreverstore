package main

import (
	"log"

	"github.com/Subham-Kedia/foreverstore/p2p"
)

func makeServer(addr string, nodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    addr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "foreverstore",
		PathTransformFunc: CASPathTransfomrFunc,
		Transport:         tr,
		bootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)
	tr.OnPeer = s.OnPeer

	return s
}

func main() {
	server1 := makeServer(":3000")
	server2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(server1.Start())
	}()

	server2.Start()
}
