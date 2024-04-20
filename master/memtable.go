package master

import "github.com/google/btree"

type memtable struct {
	data btree.BTree
}

type memtableRecord struct {
	key   []byte
	value []byte
}

func (m *memtable) Put(key, value []byte) {

}
