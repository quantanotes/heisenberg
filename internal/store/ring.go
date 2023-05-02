package store

import (
	"heisenberg/internal"
	"sort"
)

// Circular data structure used for consistent hashing
type ring struct {
	nodes  map[uint32]string
	sorted []uint32
}

func newRing() *ring {
	return &ring{
		nodes: make(map[uint32]string),
	}
}

func (r *ring) getNode(key []byte) string {
	if len(r.nodes) == 0 {
		return ""
	}

	h := internal.Hash(key)
	i := sort.Search(len(r.sorted), func(i int) bool {
		return r.sorted[i] >= h
	})

	if i == len(r.sorted) {
		i = 0
	}

	return r.nodes[r.sorted[i]]
}

func (r *ring) addNode(node string) {
	r.nodes[internal.Hash([]byte(node))] = node
	r.updatedSorted()
}

func (r *ring) removeNode(node string) {
	delete(r.nodes, internal.Hash([]byte(node)))
	r.updatedSorted()
}

func (r *ring) updatedSorted() {
	r.sorted = make([]uint32, 0, len(r.nodes))
	for k := range r.nodes {
		r.sorted = append(r.sorted, k)
	}
	sort.Slice(r.sorted, func(i, j int) bool {
		return r.sorted[i] < r.sorted[j]
	})
}
