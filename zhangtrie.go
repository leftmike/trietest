package trietest

import (
	// Clone of github.com/zhangchiqing/merkle-patricia-trie
	"github.com/leftmike/merklepatriciatrie"
)

type zhangTrie struct {
	trie *merklepatriciatrie.Trie
}

func NewZhangTrie() Trie {
	return zhangTrie{
		trie: merklepatriciatrie.NewTrie(),
	}
}

func (_ zhangTrie) Delete(key []byte) error {
	return ErrNotSupported
}

func (zt zhangTrie) Get(key []byte) ([]byte, error) {
	val, found := zt.trie.Get(key)
	if !found {
		return nil, ErrNotFound
	}
	return val, nil
}

func (zt zhangTrie) Hash() []byte {
	return zt.trie.Hash()
}

func (zt zhangTrie) Put(key, val []byte) error {
	zt.trie.Put(key, val)
	return nil
}

func (zt zhangTrie) Serialize() ([]byte, bool) {
	return zt.trie.Root().Serialize(), true
}
