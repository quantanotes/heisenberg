package master

import (
	"bytes"

	"github.com/google/btree"
)

type memtable struct {
	data         btree.BTree
	rotationData btree.BTree
	rotating     bool
}

type memtableRecord struct {
	key, value   []byte
	page, offset uint32
}

func (m *memtable) Put(key, value []byte) {
	m.data.ReplaceOrInsert(memtableRecord{key: key, value: value})
}

func (m *memtable) UpdateHeapPointer(key []byte, page, offset uint32) {
	mr := m.data.Get(memtableRecord{key: key}).(memtableRecord)
	m.data.ReplaceOrInsert(memtableRecord{key: key, value: mr.value, page: page, offset: offset})
}

func (m *memtable) Rotate() {
	m.rotating = true
	m.data, m.rotationData = m.rotationData, m.data
	m.data = *btree.New(2)
}

func (m *memtable) Foreach() {
}

func (mr memtableRecord) Less(than btree.Item) bool {
	switch than := than.(type) {
	case *memtableRecord:
		return bytes.Compare(mr.key, than.key) < 1
	default:
		return true
	}
}
