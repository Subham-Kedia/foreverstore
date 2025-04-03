package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/Subham-Kedia/foreverstore/message"
	"github.com/Subham-Kedia/foreverstore/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc TransformPathFunc
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
	// Register DataMessage type with gob
	gob.Register(&message.DataMessage{})

	storeOpts := StoreOpts{
		Root:          opts.StorageRoot,
		TransformPath: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitCh:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) broadcast(p *message.Message) error {
	for _, peer := range s.peers {
		if peer.IsOutbound() {
			if err := gob.NewEncoder(peer).Encode(p); err != nil {
				fmt.Println("Error encoding message:", err)
				return err
			}
		}
	}
	return nil
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

	p := &message.DataMessage{
		Key:  key,
		Data: buf.Bytes(),
	}
	fmt.Println(buf.Bytes())
	return s.broadcast(&message.Message{
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

// bootstrap the network by dialing the bootstrap nodes
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

// consume messages from the transport
func (s *FileServer) loop() {
	defer func() {
		fmt.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Printf("Recieved message from %v\n", msg.Payload)
			// var m message.Message
			// if len(msg.Payload) == 0 {
			// 	log.Println("Error decoding message: empty payload")
			// 	continue
			// }
			// if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&m); err != nil {
			// 	log.Println("Error decoding message:", err)
			// 	continue
			// }
			if err := s.handleMessage(&msg); err != nil {
				log.Println(err)
			}
		case <-s.quitCh:
			return
		}
	}
}

func (s *FileServer) handleMessage(m *message.Message) error {
	switch v := m.Payload.(type) {
	case *message.DataMessage:
		fmt.Printf("Data recieved %v\n", v)
		data, ok := m.Payload.(*message.DataMessage)
		if !ok {
			return fmt.Errorf("invalid payload type: %T", m.Payload)
		}
		return s.StoreData(data.Key, bytes.NewReader(data.Data))
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
