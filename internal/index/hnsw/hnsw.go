//go:generate go tool cgo hnsw.go
package hnsw

//#cgo CFLAGS: -I./ -g -Wall
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

type hnsw struct {
	index     	C.HNSW
	spaceType 	internal.spaceType
	dim       	int
	max       	int
	normalise 	bool
	opts      	hnswOptions
}

func newHNSWOptions(m int, ef int) *hnswOptions {
	return &hnswOptions{
		m:  m,
		ef: ef,
	}
}

func newHNSW(space internal.spaceType, dim int, max int, opts *hnswOptions, seed int) *hnsw {
	return &hnsw{
		index:     C.initHNSW(dim, max, opts.m, opts.ef, seed),
		normalise: space == index.cosine,
		opts:      *opts,
	}
}

func (h *hnsw) loadHNSW (path string) error {
	ret := C.loadHNSW(C.CString(path), h.dim, h.space)
	if ret == nil {
		return HNSWOperationError("loadHNSW failed")
	}
	return nil
}

func (h *hnsw) saveHNSW (path string) error {
	ret := C.saveHNSW(h.index, C.CString(path))
	if ret == false {
		return HNSWOperationError("saveHNSW failed")
	}
	return nil
}

func (h *hnsw) add(id int, vec []float32) error {
	C.addPoint(h.index, (*C.float)(unsafe.Pointer(&vec[0])), C.ulong(id))
	return nil
}

func (h *HNSW) delete(id int) error {
	C.deletePoint(h.index, C.ulong(id))
	return nil
}

func (h *hnsw) search(query []float32, k int) ([]uint32, error) {
	labels = make([]int, k)
	len = C.search(h.index, (*C.float)(unsafe.Pointer(&query[0])), k, (*C.int)(unsafe.Pointer(&labels[0])))
	if len == -1 {
		return nil, HNSWOperationError("search failed")
	}
	return labels, nil
}
