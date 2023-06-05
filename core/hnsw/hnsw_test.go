package hnsw

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/quantanotes/heisenberg/common"
	"github.com/quantanotes/heisenberg/math"
)

func TestHNSW(t *testing.T) {
	n := 1000
	opts := HNSWOptions{24, 900, n}
	config := common.IndexConfig{Name: "", Indexer: 0, FreeList: make([]uint, 0), Dim: 3, Space: int(common.Cosine), Count: 0}
	h := New(common.L2, config, opts)
	vecs := make([][]float32, n)

	for i := 0; i < n-1; i++ {
		vec := []float32{rand.Float32() - 0.5, rand.Float32() - 0.5, rand.Float32() - 0.5}
		vecs[i] = vec
		h.Insert(uint(i), vec)
	}

	query := []float32{0.1, -0.2, 0.4}
	result, err := h.Search(query, 10)
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range result {
		fmt.Printf("%v, %v\n", vecs[i], math.L2(vecs[i], query))
	}
}
