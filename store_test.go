package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestHas(t *testing.T) {
	opts := StoreOpts{
		TransformPath: CASPathTransfomrFunc,
	}
	s := NewStore(opts)
	key := "testfile"
	data := []byte("This is a testfile data")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Errorf("write failed: %s", err)
	}

	err := s.Has(key)
	if err {
		t.Errorf("has failed")
	}
	// s.Clear()
}
func TestRead(t *testing.T) {
	opts := StoreOpts{
		TransformPath: CASPathTransfomrFunc,
	}
	s := NewStore(opts)
	key := "testfile"
	data := []byte("This is a testfile data")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Errorf("write failed: %s", err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Errorf("read failed: %s", err)
	}
	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("read data mismatch: got %s, want %s", string(b), string(data))
	}
}
func TestWrite(t *testing.T) {
	opts := StoreOpts{
		TransformPath: CASPathTransfomrFunc,
	}
	s := NewStore(opts)
	key := "testfile1"
	data := []byte("This is a testfile data")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Errorf("write failed, %s", err)
	}
}
func TestClear(t *testing.T) {
	opts := StoreOpts{
		TransformPath: CASPathTransfomrFunc,
	}
	s := NewStore(opts)
	key := "bkl"
	data := []byte("This is a test file")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error("write failed")
	}

	if err := s.Clear(); err != nil {
		t.Errorf("clear failed, %s", err)
	}

}
func TestDeleteFile(t *testing.T) {
	opts := StoreOpts{
		TransformPath: CASPathTransfomrFunc,
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
		TransformPath: CASPathTransfomrFunc,
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
