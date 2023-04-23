package main

import (
	"heisenberg/hnsw"
	"math/rand"
)

type indexConfig struct {
	dim      int    // Dimension size of vectors
	maxSize  int    // Max amount of vectors in index
	space    string // Distance measure
	cursor   int    // Position to insert vector
	freeList []int  // Queue data structure to store deleted vectors
}

type Index struct {
	path   string       // Filepath to save/load index from
	hnsw   *hnsw.HNSW   // C++ implementation of hnsw indexing algorithm
	config *indexConfig // On-kv configuration
}

func NewIndex(path string, config *indexConfig, m int, ef int) *Index {
	hnsw := hnsw.New(config.dim, m, ef, rand.Int(), uint32(config.maxSize), config.space)
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
	hnsw := hnsw.Load(path, config.dim, config.space)
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
	idx.hnsw.Free()
	idx = nil
}

func (idx *Index) Save() {
	idx.hnsw.Save(idx.path)
}

func (idx *Index) Put(vec []float32, id uint32) {
	idx.hnsw.AddPoint(vec, id)
}

// Get next available index position without mutating index state
func (idx *Index) Next() int {
	if len(idx.config.freeList) > 0 {
		id := idx.config.freeList[0]
		return id
	}

	id := idx.config.cursor
	return id
}

// Get next available index position and mutate index state
func (idx *Index) MutNext() int {
	if len(idx.config.freeList) > 0 {
		id := idx.config.freeList[0]
		idx.config.freeList = idx.config.freeList[0:]
		return id
	}

	id := idx.config.cursor
	idx.config.cursor++
	return id
}
