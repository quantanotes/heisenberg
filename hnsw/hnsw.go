//go:generate go tool cgo hnsw.go
package hnsw

//#cgo CFLAGS: -I./ -g -Wall
//#cgo CXXFLAGS: -std=c++11
//#include <stdlib.h>
//#include "hnsw_wrapper.h"
import "C"

import (
	"fmt"
	"github.com/quantanotes/heisenberg/math"
	"github.com/quantanotes/heisenberg/utils"
	"unsafe"
)

type HNSWOptions struct {
	M  int
	Ef int
}

type HNSW struct {
	index     C.HNSW
	space     utils.SpaceType
	dim       int
	normalise bool
}

func NewHNSW(space utils.SpaceType, dim int, max int, opts *HNSWOptions, seed int) *HNSW {
	return &HNSW{
		index:     C.initHNSW(C.int(dim), C.ulong(max), C.int(opts.M), C.int(opts.Ef), C.int(seed), C.int(space)),
		space:     space,
		dim:       dim,
		normalise: space == utils.Cosine,
	}
}

func LoadHNSW(path string, dim int, space utils.SpaceType) (*HNSW, error) {
	index := C.loadHNSW(C.CString(path), C.int(dim), C.int(space))
	if index == nil {
		return nil, fmt.Errorf("hnsw load failed")
	}
	normalise := space == utils.Cosine
	return &HNSW{
		index,
		space,
		dim,
		normalise,
	}, nil
}

func (h *HNSW) Free() {
	C.freeHNSW(h.index)
}

func (h *HNSW) Save(path string) error {
	res := C.saveHNSW(h.index, C.CString(path))
	if !res {
		return fmt.Errorf("hnsw save failed")
	}
	return nil
}

func (h *HNSW) Add(id int, vec []float32) error {
	if h.normalise {
		vec = math.Normalise(vec)
	}
	C.addPoint(h.index, (*C.float)(unsafe.Pointer(&vec[0])), C.ulong(id))
	return nil
}

func (h *HNSW) Delete(id int) error {
	C.deletePoint(h.index, C.ulong(id))
	return nil
}

func (h *HNSW) Search(query []float32, k int) ([]int, error) {
	cids := make([]C.ulong, k)
	cdists := make([]C.float, k)
	len := int(C.search(h.index, (*C.float)(unsafe.Pointer(&query[0])), C.int(k), &cids[0], &cdists[0]))
	if h.normalise {
		query = math.Normalise(query)
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
