package hnsw

import (
	"fmt"
	"heisenberg/internal"
	"heisenberg/internal/vector"
	"math/rand"
	"testing"
)

func TestHNSW(t *testing.T) {
	n := 1000
	opts := &hnswOptions{24, 900}
	h := newHNSW(internal.L2, 3, n, opts, 1)
	vecs := make([][]float32, n)
	for i := 0; i < n-1; i++ {
		vec := []float32{rand.Float32() - 0.5, rand.Float32() - 0.5, rand.Float32() - 0.5}
		vecs[i] = vec
		h.add(vec, i)
	}

	query := []float32{0.1, -0.2, 0.4}
	result, err := h.search(query, 10)
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range result {
		fmt.Printf("%v, %v\n", vecs[i], vector.L2(vecs[i], query))
	}
}
