package main

import (
	"bytes"
	"log"
	"time"

	"github.com/Subham-Kedia/foreverstore/p2p"
)

func makeServer(addr string, nodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    addr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
	}
	// transport
	tr := p2p.NewTCPTransport(tcpOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "foreverstore",
		PathTransformFunc: CASPathTransfomrFunc,
		Transport:         tr,
		bootstrapNodes:    nodes,
	}

	server := NewFileServer(fileServerOpts)
	tr.OnPeer = server.OnPeer

	return server
}

func main() {
	server1 := makeServer(":3000")
	server2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(server1.Start())
	}()

	go server2.Start()

	time.Sleep(time.Second * 2)

	data := bytes.NewReader([]byte("this is a test data"))
	server2.StoreData("file1", data)

  select{}
}
