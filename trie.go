package trietest

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("trietest: not found")
	ErrNotSupported = errors.New("trietest: not supported")
)

type Trie interface {
	Delete(key []byte) error
	Get(key []byte) ([]byte, error)
	Hash() []byte
	Put(key, val []byte) error
	Serialize() ([]byte, bool)
}
