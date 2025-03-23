package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
    ListenAddr: ":4000",
    HandshakeFunc: NOPHandshakeFunc,
    Decoder: NOPDecoder{},
  }
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddr, opts.ListenAddr)
	assert.Nil(t, tr.ListenAndAccept())
} 
