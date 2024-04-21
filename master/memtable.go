package master

import (
	"bytes"

	"github.com/google/btree"
)

type memtable struct {
	data         btree.BTree
	rotationData btree.BTree
	rotating     bool
	indexBuilder memtableIndexBuilder
}

type memtableEntry struct {
	key, value   []byte
	page, offset uint32
}

func (m *memtable) Put(key, value []byte) {
	m.data.ReplaceOrInsert(memtableEntry{key: key, value: value})
}

func (m *memtable) UpdateHeapPointer(key []byte, page, offset uint32) {
	me := m.data.Get(memtableEntry{key: key}).(memtableEntry)
	m.data.ReplaceOrInsert(memtableEntry{key: key, value: me.value, page: page, offset: offset})
}

func (m *memtable) Rotate() {
	m.rotating = true
	m.data, m.rotationData = m.rotationData, m.data
	m.data = *btree.New(2)
}

func (me memtableEntry) Less(than btree.Item) bool {
	switch than := than.(type) {
	case *memtableEntry:
		return bytes.Compare(me.key, than.key) < 1
	default:
		return true
	}
}
