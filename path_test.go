package main

import (
	"strings"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	root := "foreverstore"
	storageKey := "file1"

	pathKey := CASPathTransfomrFunc(storageKey, root)

	if !strings.HasPrefix(pathKey.PathName, root) {
		t.Errorf("Expected PathName to start with root '%s', got '%s'", root, pathKey.PathName)
	}

	if len(pathKey.FileName) != 40 { // SHA-1 hash length in hex
		t.Errorf("Expected FileName to be 40 characters long, got %d", len(pathKey.FileName))
	}
}

func TestCASPathTransformFuncDifferentKeys(t *testing.T) {
	root := "foreverstore"
	key1 := "file1"
	key2 := "file2"

	pathKey1 := CASPathTransfomrFunc(key1, root)
	pathKey2 := CASPathTransfomrFunc(key2, root)

	if pathKey1.FileName == pathKey2.FileName {
		t.Errorf("Expected different FileNames for different keys, got '%s' and '%s'", pathKey1.FileName, pathKey2.FileName)
	}
}

func TestCASPathTransformFuncSameKey(t *testing.T) {
	root := "foreverstore"
	key := "file1"

	pathKey1 := CASPathTransfomrFunc(key, root)
	pathKey2 := CASPathTransfomrFunc(key, root)

	if pathKey1.FileName != pathKey2.FileName {
		t.Errorf("Expected same FileName for the same key, got '%s' and '%s'", pathKey1.FileName, pathKey2.FileName)
	}
}
