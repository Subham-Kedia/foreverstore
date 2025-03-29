package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/Subham-Kedia/foreverstore/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	TCPTransportOpts  p2p.TCPTransportOpts

	bootstrapNodes []string
}

type FileServer struct {
	FileServerOpts
	store *Store

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	quitCh   chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:          opts.StorageRoot,
		PathTransform: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitCh:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

type Payload struct {
	data []byte
	key  string
}

func (s *FileServer) broadcast(p Payload) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter()
	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
  if err := s.store.writeStream(key, r); err != nil {
    return err
  }
  buf := new(bytes.Buffer)
  _, err := io.Copy(buf, r)
  
  if err != nil {
    return err
  }
  
  p := Payload{
    key: key,
    data: buf.Bytes(),
  }
  fmt.Println(buf.Bytes())
	return s.broadcast(p)
}

func (s *FileServer) Stop() {
	close(s.quitCh)
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[peer.RemoteAddr().String()] = peer
	log.Printf("connected with remote %s", peer.RemoteAddr())
	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.bootstrapNodes {
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error", addr)
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		fmt.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			// handling message from incoming connections
			fmt.Println(msg)
		case <-s.quitCh:
			return
		}
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	s.loop()
	return nil
}
