package master

import (
	"heisenberg/model"

	"github.com/google/btree"
)

// Builds a SSTable during rotation
type memtableIndexBuilder struct {
	*memtable
}

func (mib *memtableIndexBuilder) Build() model.Index {
	index := model.Index{}
	blocks := mib.chunk()

	for _, block := range blocks {
		indexBlock := model.IndexBlock{}
		for _, me := range block {
			ie := model.IndexEntry{
				Key:    me.key,
				Page:   me.page,
				Offset: me.offset,
			}
			indexBlock.Entries = append(indexBlock.Entries, ie)
		}
		index.Blocks = append(index.Blocks, indexBlock)
	}

	index.Header.First = blocks[0][0].key
	index.Header.Last = blocks[len(blocks)-1][len(blocks[len(blocks)-1])-1].key
	return index
}

func (mib *memtableIndexBuilder) chunk() [][]memtableEntry {
	var blocks [][]memtableEntry
	var currentBlock []memtableEntry

	mib.rotationData.Ascend(func(item btree.Item) bool {
		entry := item.(memtableEntry)
		currentBlock = append(currentBlock, entry)
		if len(currentBlock) == model.IndexBlockSize {
			blocks = append(blocks, currentBlock)
			currentBlock = nil
		}
		return true
	})

	if len(currentBlock) > 0 {
		blocks = append(blocks, currentBlock)
	}

	return blocks
}
