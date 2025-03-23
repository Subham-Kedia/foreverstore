package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func CASPathTransfomrFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)
	for i := range sliceLen {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Original: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Original string
}

func (p PathKey) FileName() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Original)
}

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
	pathkey := s.PathTransform(key)
	if err := os.MkdirAll(pathkey.Pathname, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(pathkey.FileName())
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

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathkey := s.PathTransform(key)
	return os.Open(pathkey.FileName())
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
