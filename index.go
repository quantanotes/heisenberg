package main

import (
	"heisenberg/hnsw"
	"math/rand"
)

type Index struct {
	path string
	hnsw *hnsw.HNSW
}

func NewIndex(path string, dim int, m int, ef int, space string) *Index {
	hnsw := hnsw.New(dim, m, ef, rand.Int(), 1000000, space)

	return &Index{
		path: path,
		hnsw: hnsw,
	}
}

func LoadIndex(path string, dim int, space string) *Index {
	hnsw := hnsw.Load(path, dim, space)
	return &Index{
		path: path,
		hnsw: hnsw,
	}
}

func (i *Index) Search(vec []float32, k int) []uint32 {
	ids, _ := i.hnsw.SearchKNN(vec, k)
	return ids
}

func (i *Index) Put(vec []float32, id int) {
	i.hnsw.AddPoint(vec, uint32(id))
}

func (i *Index) Delete() {

}
