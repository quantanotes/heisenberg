//go:generate go tool cgo hnsw.go
package hnsw

//#cgo CFLAGS: -I./
//#cgo LDFLAGS: -lhnsw
//#include <stdlib.h>
//#include "hnsw_wrapper.h"
import "C"

import (
	"heisenberg/internal"
	"unsafe"
)

type hnswOptions struct {
	m  int
	ef int
}

type HNSW struct {
	index     C.HNSW
	space     internal.SpaceType
	dim       int
	max       int
	normalise bool
	opts      hnswOptions
}

func newHnsw(space internal.SpaceType, dim int, max int, opts *hnswOptions) *HNSW {
	return &HNSW{
		index:     C.new_hnsw(),
		normalise: space == internal.Cosine,
		opts:      *opts,
	}
}

func (h *HNSW) init(space internal.SpaceType, dim int, max int, opts *hnswOptions) {
	hnsw := newHnsw(space, dim, max, opts)
	C.init_hnsw(hnsw.index, C.int(dim), C.int(max), C.bool(hnsw.normalise), &hnsw.opts.m, &hnsw.opts.ef, c)
}

func (h *HNSW) add(id int, vec []float32) error {
	C.addPoint(h.index, (*C.float)(unsafe.Pointer(&vec[0])), C.ulong(id))
	return nil
}

func (h *HNSW) delete(id int) error {
	C.deletePoint(h.index, C.ulong(id))
	return nil
}

func (h *HNSW) search(query []float32, k int) ([]int, error) {
	C.search(h.index, (*C.float)(unsafe.Pointer(&query[0])), C.ulong(k))
	return nil, nil
}
