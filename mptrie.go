package trietest

import (
	"github.com/leftmike/mptrie"
)

type mpTrie struct {
	trie *mptrie.MPTrie
}

func NewMPTrie() Trie {
	return mpTrie{
		trie: mptrie.New(),
	}
}

func (mpt mpTrie) Delete(key []byte) error {
	err := mpt.trie.Delete(key)
	if err == mptrie.ErrNotFound {
		return ErrNotFound
	}
	return err
}

func (mpt mpTrie) Get(key []byte) ([]byte, error) {
	val, err := mpt.trie.Get(key)
	if err == mptrie.ErrNotFound {
		return nil, ErrNotFound
	}
	return val, err
}

func (mpt mpTrie) Put(key, val []byte) error {
	return mpt.trie.Put(key, val)
}
