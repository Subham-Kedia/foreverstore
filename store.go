package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type StoreOpts struct {
	Root          string
	PathTransform PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	fullPath := s.PathTransform(key).FullPath()

	_, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return true
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathkey := s.PathTransform(key)
	return os.Open(pathkey.FullPath())
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathkey := s.PathTransform(key)
	if err := os.MkdirAll(pathkey.PathName, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(pathkey.FullPath())
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	fmt.Printf("wrote %d bytes\n", n)
	return nil
}

func (s *Store) Delete(key string) error {
	pathkey := s.PathTransform(key)
	defer func() {
		log.Printf("deleted %s from disk", pathkey.PathName)
	}()
	return os.RemoveAll(pathkey.PathName)
}

func (s *Store) Clear() error {
  return os.RemoveAll(s.Root)
}
