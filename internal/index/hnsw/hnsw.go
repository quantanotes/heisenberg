//go:generate go tool cgo hnsw.go
package hnsw

//#cgo CFLAGS: -I./ -g -Wall
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
	index     	C.HNSW
	spaceType 	int
	dim       	int
	max       	int
	normalise 	bool
	opts      	hnswOptions
}

func newHNSW(space int, dim int, max int, opts *hnswOptions) *hnsw {
	return &hnsw{
		index:     C.initHNSW(),
		normalise: space == index.cosine,
		opts:      *opts,
	}
}

func (h *hnsw) loadHNSW (path string) {
	C.loadHNSW(C.CString(path), h.dim, h.space)
}

func (h *hnsw) addPoint(id int, vec []float32) error {
	C.addPoint(h.index, (*C.float)(unsafe.Pointer(&vec[0])), C.ulong(id))
	return nil
}

func (h *hnsw) deletePoint(id int) error {
	C.deletePoint(h.index, C.ulong(id))
	return nil
}

func (h *hnsw) search(query []float32, k int) ([]uint32, error) {
	labels = make([]uint32, k)
	len = C.search(h.index, (*C.float)(unsafe.Pointer(&query[0])), k, (*C.int)(unsafe.Pointer(&labels[0])))
	return labels, nil
}
