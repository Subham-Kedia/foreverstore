package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransform: DefaultPathTransformFunc,
	}
	store := NewStore(opts)
	if err := store.writeStream("bucket2", bytes.NewReader([]byte("This is a distributed file system"))); err != nil {
		t.Fatalf("error writing stream: %s", err)
	}
}
