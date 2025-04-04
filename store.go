package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

type StoreOpts struct {
	Root          string
	TransformPath TransformPathFunc
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
	fullPath := s.TransformPath(key, s.Root).FullPath()

	_, err := os.Stat(fullPath)
	return err == nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathkey := s.TransformPath(key, s.Root)
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
	pathkey := s.TransformPath(key, s.Root)
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
	log.Printf("wrote %d bytes\n", n)
	return nil
}

func (s *Store) Delete(key string) error {
	pathkey := s.TransformPath(key, s.Root)
	defer func() {
		log.Printf("deleted %s from disk", pathkey.PathName)
	}()
	return os.RemoveAll(pathkey.PathName)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}
