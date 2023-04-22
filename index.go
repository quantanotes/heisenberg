package main

import "heisenberg/hnsw"

type Index struct {
	hnsw hnsw.HNSW
}

func NewIndex() *Index {
	hnsw := hnsw.New(1, 1, 1, 1, 1, "cosine")

	return &Index{
		hnsw: hnsw,
	}
}
