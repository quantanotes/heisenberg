package main

import "heisenberg/hnsw"

type Index struct {
	path string
	hnsw *hnsw.HNSW
}
