package main

import (
	"bytes"
	"fmt"
	"io"
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

	tr := p2p.NewTCPTransport(tcpOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       addr[1:] + "_foreverstore",
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
	server3 := makeServer(":5001")
	server2 := makeServer(":4000", ":3000", ":5001")

	go func() {
		log.Fatal(server1.Start())
	}()

	go func() {
		log.Fatal(server3.Start())
	}()

	time.Sleep(time.Second * 3)

	go func() {
		log.Fatal(server2.Start())
	}()

	time.Sleep(time.Second * 2)

	data := bytes.NewReader([]byte("this is a confidential data"))
	server2.StoreData("testfile", data)

	time.Sleep(time.Second * 2)

	r, _ := server3.store.Read("testfile")
	b, _ := io.ReadAll(r)
	fmt.Println(string(b))

	select {}
}
