package main

import (
	"fmt"
	"io"
	"os"
)

type PathTransformFunc func(string) string
type StoreOpts struct {
	PathTransform PathTransformFunc
}

type Store struct {
	StoreOpts
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathname := s.PathTransform(key)
	if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
		return err
	}
	filename := "somefilename"
	f, err := os.Create(pathname + "/" + filename)
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

func (s *Store) readStream(key string, w io.Writer) error {
	return nil
}
