//go:generate go tool cgo hnsw.go
package hnsw

//#cgo CFLAGS: -I./ -g -Wall
//#include <stdlib.h>
//#include "hnsw_wrapper.h"
import "C"

import (
	"fmt"
	"heisenberg/internal"
	"heisenberg/internal/vector"
	"unsafe"
)

type hnswOptions struct {
	m  int
	ef int
}

type hnsw struct {
	index     C.HNSW
	space     internal.SpaceType
	dim       int
	max       int
	normalise bool
	opts      hnswOptions
}

func newHNSW(space internal.SpaceType, dim int, max int, opts *hnswOptions, seed int) *hnsw {
	return &hnsw{
		index:     C.initHNSW(C.int(dim), C.ulong(max), C.int(opts.m), C.int(opts.ef), C.int(seed), C.int(space)),
		space:     space,
		dim:       dim,
		max:       max,
		normalise: space == internal.Cosine,
		opts:      *opts,
	}
}

func (h *hnsw) load(path string) error {
	res := C.loadHNSW(C.CString(path), C.int(h.dim), C.int(h.space))
	if res == nil {
		return fmt.Errorf("hnsw load failed")
	}
	return nil
}

func (h *hnsw) save(path string) error {
	res := C.saveHNSW(h.index, C.CString(path))
	if !res {
		return fmt.Errorf("hnsw save failed")
	}
	return nil
}

func (h *hnsw) add(vec []float32, id int) error {
	if h.normalise {
		vec = vector.Normalise(vec)
	}
	C.addPoint(h.index, (*C.float)(unsafe.Pointer(&vec[0])), C.ulong(id))
	return nil
}

func (h *hnsw) delete(id int) error {
	C.deletePoint(h.index, C.ulong(id))
	return nil
}

func (h *hnsw) search(query []float32, k int) ([]int, error) {
	cids := make([]C.ulong, k)
	len := int(C.search(h.index, (*C.float)(unsafe.Pointer(&query[0])), C.int(k), &cids[0], &cdists[0]))
	if h.normalise {
		query = vector.Normalise(query)
	}
	if len == -1 {
		return nil, fmt.Errorf("hnsw search failed")
	}
	ids := make([]int, len)
	for i := 0; i < len; i++ {
		ids[i] = int(cids[i])
	}
	return ids, nil
}
