package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTrasnformFunc(t *testing.T) {
  key := "random"  
  pathname := CASPathTransfomrFunc(key)
  fmt.Println(pathname)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransform: CASPathTransfomrFunc,
	}
	store := NewStore(opts)
  key := "specialkey"
  data := []byte("This is a text to be in this file")
	if err := store.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Fatalf("error writing stream: %s", err)
	}

  r, err := store.Read(key)
  if err != nil {
    t.Error(err)
  }
  b, _ := io.ReadAll(r)
  fmt.Println(b)
}
