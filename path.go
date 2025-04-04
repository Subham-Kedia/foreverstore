package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

type PathKey struct {
	PathName string
	FileName string
}

type TransformPathFunc func(string, string) PathKey

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func CASPathTransfomrFunc(key string, root string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 8
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)
	for i := range sliceLen {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: fmt.Sprintf("%s/%s", root, strings.Join(paths, "/")),
		FileName: hashStr,
	}
}
