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

type Message struct {
	From    string
	Payload any
}

type DataMessage struct {
	data []byte
	key  string
}

func (s *FileServer) broadcast(p *Message) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)
	if err := s.store.writeStream(key, tee); err != nil {
		return err
	}
	_, err := io.Copy(buf, r)

	if err != nil {
		return err
	}

	p := &DataMessage{
		key:  key,
		data: buf.Bytes(),
	}
	fmt.Println(buf.Bytes())
	return s.broadcast(&Message{
		From:    "todo",
		Payload: p,
	})
}

func (s *FileServer) Stop() {
	close(s.quitCh)
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[peer.RemoteAddr().String()] = peer
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
			var m Message
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&m); err != nil {
				log.Fatal(err)
			}

			if err := s.handleMessage(&m); err != nil {
				log.Println(err)
			}
		case <-s.quitCh:
			return
		}
	}
}

func (s *FileServer) handleMessage(msg *Message) error {
  switch v := msg.Payload.(type) {
    case *DataMessage:
      fmt.Printf("Data recieved %v\n", v)
  }
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	s.loop()
	return nil
}
