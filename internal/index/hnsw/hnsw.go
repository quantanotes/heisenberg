//go:generate go tool cgo hnsw.go
package hnsw

//#cgo CFLAGS: -I./
//#cgo LDFLAGS: -lhnsw
//#include <stdlib.h>
//#include "hnsw_wrapper.h"
import "C"
import "unsafe"

type hnswOptions struct {
	m  int
	ef int
}

type hnsw struct {
	index     C.HNSW
	space     spaceType
	dim       int
	max       int
	normalise bool
	opts      hnswOptions
}

func newHnsw(space spaceType, dim int, max int, opts *hnswOptions) *hnsw {
	return &hnsw{
		index:     C.new_hnsw(),
		normalise: space == index.cosine,
		opts:      *opts,
	}
}

func (h *hnsw) add(id int, vec []float32) error {
	C.addPoint(h.index, (*C.float)(unsafe.Pointer(&vec[0])), C.ulong(id))
	return nil
}

func (h *hnsw) delete(id int) error {
	C.deletePoint(h.index, C.ulong(id))
	return nil
}

func (h *hnsw) search(query []float32, k int) ([]int, error) {
	C.search(h.index, (*C.float)(unsafe.Pointer(&query[0])), C.ulong(k))
	return nil, nil
}
