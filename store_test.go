package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTrasnformFunc(t *testing.T) {
	key1 := "ab"
	key2 := "abc"
	key3 := "abcd"
	pathname1 := CASPathTransfomrFunc(key1)
	fmt.Println(pathname1)
	pathname2 := CASPathTransfomrFunc(key2)
	fmt.Println(pathname2)
	pathname3 := CASPathTransfomrFunc(key3)
	fmt.Println(pathname3)
}

func TestDeleteFile(t *testing.T) {
	opts := StoreOpts{
		PathTransform: CASPathTransfomrFunc,
	}
	s := NewStore(opts)
	key := "bkl"
	data := []byte("This is a test file")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error("write failed")
	}

	if err := s.Delete(key); err != nil {
		t.Errorf("delete failed, %s", err)
	}

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
