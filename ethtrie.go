package trietest

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	ethtrie "github.com/ethereum/go-ethereum/trie"
)

type ethTrie struct {
	trie *ethtrie.Trie
}

func NewEthTrie() Trie {
	trie, err := ethtrie.New(common.Hash{}, ethtrie.NewDatabase(memorydb.New()))
	if err != nil {
		panic(fmt.Sprintf("ethtrie: %s", err))
	}

	return ethTrie{
		trie: trie,
	}
}

func (et ethTrie) Delete(key []byte) error {
	val, err := et.trie.TryGet(key)
	if len(val) == 0 {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	return et.trie.TryDelete(key)
}

func (et ethTrie) Get(key []byte) ([]byte, error) {
	val, err := et.trie.TryGet(key)
	if len(val) == 0 {
		return nil, ErrNotFound
	}
	return val, err
}

func (et ethTrie) Hash() []byte {
	h := et.trie.Hash()
	return h[:]
}

func (et ethTrie) Put(key, val []byte) error {
	return et.trie.TryUpdate(key, val)
}

func (_ ethTrie) Serialize() ([]byte, bool) {
	return nil, false
}
