package main

import (
	"heisenberg/hnsw"
	"math/rand"
)

type indexConfig struct {
	Dim      int    `json:"dim"`      // Dimension size of vectors
	MaxSize  int    `json:"maxSize"`  // Max amount of vectors in index
	Space    string `json:"space"`    // Distance measure
	Cursor   int    `json:"cursor"`   // Position to insert vector
	FreeList []int  `json:"freeList"` // Queue data structure to store deleted vectors
}

type Index struct {
	path   string       // Filepath to save/load index from
	hnsw   *hnsw.HNSW   // C++ implementation of hnsw indexing algorithm
	config *indexConfig // On-kv configuration
}

func NewIndex(path string, config *indexConfig, m int, ef int) *Index {
	hnsw := hnsw.New(config.Dim, m, ef, rand.Int(), uint32(config.MaxSize), config.Space)
	if hnsw == nil {
		return nil
	}

	return &Index{
		path,
		hnsw,
		config,
	}
}

func LoadIndex(path string, config *indexConfig) *Index {
	hnsw := hnsw.Load(path, config.Dim, config.Space)
	if hnsw == nil {
		return nil
	}

	return &Index{
		path,
		hnsw,
		config,
	}
}

func (idx *Index) Close() {
	if idx.hnsw != nil {
		idx.hnsw.Free()
	}
	idx = nil
}

func (idx *Index) Save() {
	if idx.hnsw != nil {
		idx.hnsw.Save(idx.path)
	}
}

func (idx *Index) Put(vec []float32, id uint32) {
	idx.hnsw.AddPoint(vec, id)
}

func (idx *Index) Search(vec []float32, k int) []uint32 {
	ids, _ := idx.hnsw.SearchKNN(vec, k)
	return ids
}

// Get next available index position without mutating index state
func (idx *Index) Next() int {
	if len(idx.config.FreeList) > 0 {
		id := idx.config.FreeList[0]
		return id
	}

	id := idx.config.Cursor
	return id
}

// Get next available index position and mutate index state
func (idx *Index) MutNext() int {
	if len(idx.config.FreeList) > 0 {
		id := idx.config.FreeList[0]
		idx.config.FreeList = idx.config.FreeList[0:]
		return id
	}

	id := idx.config.Cursor
	idx.config.Cursor++
	return id
}

func (idx *Index) Delete(id int) {
	idx.config.FreeList = append(idx.config.FreeList, id)
}
