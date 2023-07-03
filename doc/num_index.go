package doc

import (
	"fmt"
	"sort"
)

type indexEntry struct {
	id    uint64
	value float64
}

type numIndex interface {
	insert(id uint64, value float64) error
	remove(id uint64) error
	search(comp comparison, value float64, limit uint) []uint64
}

type diskNumIndex struct {
}

type memNumIndex struct {
	index []indexEntry
}

func newMemNumIndex() numIndex {
	return &memNumIndex{make([]indexEntry, 0)}
}

func (idx *memNumIndex) insert(id uint64, value float64) error {
	insertIdx := sort.Search(len(idx.index), func(i int) bool {
		return idx.index[i].value > value
	})
	idx.index = append(idx.index[:insertIdx+1], idx.index[insertIdx:]...)
	idx.index[insertIdx] = indexEntry{id, value}
	return nil
}

func (idx *memNumIndex) remove(id uint64) error {
	removeIdx := sort.Search(len(idx.index), func(i int) bool {
		return idx.index[i].id == id
	})
	if removeIdx == -1 {
		return fmt.Errorf("id %v not in index", id)
	}
	idx.index = append(idx.index[:removeIdx], idx.index[removeIdx+1:]...)
	return nil
}

func (idx *memNumIndex) search(comp comparison, value float64, limit uint) []uint64 {
	result := make([]uint64, 0, limit)

	searchIdx := sort.Search(len(idx.index), func(i int) bool {
		return comp(idx.index[i].value, value)
	})

	for i := searchIdx; i < len(idx.index); i++ {
		if len(result) >= int(limit) {
			break
		}
		if comp(idx.index[i].value, value) {
			result = append(result, idx.index[i].id)
		}
	}

	return result
}
