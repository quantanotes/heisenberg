//go:generate go tool cgo hnsw.go
package hnsw

//#cgo CFLAGS: -I./ -g -Wall
//#cgo CXXFLAGS: -std=c++11
//#include <stdlib.h>
//#include "hnsw_wrapper.h"
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/quantanotes/heisenberg/common"
	"github.com/quantanotes/heisenberg/math"
)

type HNSWOptions struct {
	M   int
	Ef  int
	Max int
}

type HNSW struct {
	index     C.HNSW
	config    common.IndexConfig
	normalise bool
}

func New(space common.SpaceType, config common.IndexConfig, opts HNSWOptions) *HNSW {
	h := &HNSW{
		index:     C.initHNSW(C.int(config.Dim), C.ulong(opts.Max), C.int(opts.M), C.int(opts.Ef), C.int(42069), C.int(config.Space)),
		config:    config,
		normalise: common.SpaceType(config.Space) == common.Cosine,
	}
	return h
}

func Load(path string, config common.IndexConfig) (*HNSW, error) {
	index := C.loadHNSW(C.CString(path), C.int(config.Dim), C.int(config.Space))
	if index == nil {
		return nil, fmt.Errorf("hnsw load failed")
	}
	normalise := common.SpaceType(config.Space) == common.Cosine
	h := &HNSW{
		index,
		config,
		normalise,
	}
	return h, nil
}

func (h *HNSW) Close() {
	C.freeHNSW(h.index)
}

func (h *HNSW) Save(path string) error {
	res := C.saveHNSW(h.index, C.CString(path))
	if !res {
		return fmt.Errorf("hnsw save failed")
	}
	return nil
}

func (h *HNSW) Insert(id uint, vec []float32) error {
	if len(vec) != int(h.config.Dim) {
		return common.DimensionMismatch(len(vec), int(h.config.Dim))
	}
	if h.normalise {
		vec = math.Normalise(vec)
	}
	C.addPoint(h.index, (*C.float)(unsafe.Pointer(&vec[0])), C.ulong(id))
	return nil
}

func (h *HNSW) Delete(id uint) error {
	C.deletePoint(h.index, C.ulong(id))
	h.config.FreeList = append(h.config.FreeList, id)
	return nil
}

func (h *HNSW) Search(query []float32, k uint) ([]uint, error) {
	cids := make([]C.ulong, k)
	cdists := make([]C.float, k)
	len := int(C.search(h.index, (*C.float)(unsafe.Pointer(&query[0])), C.int(k), &cids[0], &cdists[0]))
	if h.normalise {
		query = math.Normalise(query)
	}
	if len == -1 {
		return nil, fmt.Errorf("hnsw search failed")
	}
	ids := make([]uint, len)
	for i := 0; i < len; i++ {
		ids[i] = uint(cids[i])
	}
	return ids, nil
}

func (h *HNSW) nextMut() uint {
	size := len(h.config.FreeList)
	if size > 0 {
		next := h.config.FreeList[0]
		h.config.FreeList = h.config.FreeList[1:size]
		return next
	}
	h.config.Count++
	return h.config.Count
}

func (h *HNSW) Next() uint {
	size := len(h.config.FreeList)
	if size > 0 {
		next := h.config.FreeList[0]
		h.config.FreeList = h.config.FreeList[1:size]
		return next
	}
	return h.config.Count
}

func (h *HNSW) GetConfig() common.IndexConfig {
	return h.config
}
