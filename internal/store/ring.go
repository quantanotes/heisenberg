package store

import (
	"heisenberg/internal"
	"sort"
)

// Circular data structure used for consistent hashing
type ring struct {
	Nodes  map[uint32]string
	Sorted []uint32
}

func newRing() ring {
	return ring{
		Nodes: make(map[uint32]string),
	}
}

func (r *ring) getNode(key []byte) string {
	if len(r.Nodes) == 0 {
		return ""
	}

	h := internal.Hash(key)
	i := sort.Search(len(r.Sorted), func(i int) bool {
		return r.Sorted[i] >= h
	})

	if i == len(r.Sorted) {
		i = 0
	}

	return r.Nodes[r.Sorted[i]]
}

func (r *ring) addNode(node string) {
	r.Nodes[internal.Hash([]byte(node))] = node
	r.updatedSorted()
}

func (r *ring) removeNode(node string) {
	delete(r.Nodes, internal.Hash([]byte(node)))
	r.updatedSorted()
}

func (r *ring) updatedSorted() {
	r.Sorted = make([]uint32, 0, len(r.Nodes))
	for k := range r.Nodes {
		r.Sorted = append(r.Sorted, k)
	}
	sort.Slice(r.Sorted, func(i, j int) bool {
		return r.Sorted[i] < r.Sorted[j]
	})
}
