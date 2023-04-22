package hnsw

// #cgo CFLAGS: -I./
// #cgo LDFLAGS: -lhnsw
// #include <stdlib.h>
// #include "hnsw_wrapper.h"
import "C"
import (
	"math"
	"unsafe"
)

type HNSW struct {
	index     C.HNSW
	spaceType string
	dim       int
	normalize bool
}

func New(dim, M, efConstruction, randSeed int, maxElements uint32, spaceType string) *HNSW {
	var hnsw HNSW
	hnsw.dim = dim
	hnsw.spaceType = spaceType
	if spaceType == "ip" {
		hnsw.index = C.initHNSW(C.int(dim), C.ulong(maxElements), C.int(M), C.int(efConstruction), C.int(randSeed), C.char('i'))
	} else if spaceType == "cosine" {
		hnsw.normalize = true
		hnsw.index = C.initHNSW(C.int(dim), C.ulong(maxElements), C.int(M), C.int(efConstruction), C.int(randSeed), C.char('i'))
	} else {
		hnsw.index = C.initHNSW(C.int(dim), C.ulong(maxElements), C.int(M), C.int(efConstruction), C.int(randSeed), C.char('l'))
	}
	return &hnsw
}

func Load(location string, dim int, spaceType string) *HNSW {
	var hnsw HNSW
	hnsw.dim = dim
	hnsw.spaceType = spaceType

	pLocation := C.CString(location)
	if spaceType == "ip" {
		hnsw.index = C.loadHNSW(pLocation, C.int(dim), C.char('i'))
	} else if spaceType == "cosine" {
		hnsw.normalize = true
		hnsw.index = C.loadHNSW(pLocation, C.int(dim), C.char('i'))
	} else {
		hnsw.index = C.loadHNSW(pLocation, C.int(dim), C.char('l'))
	}
	C.free(unsafe.Pointer(pLocation))
	return &hnsw
}

func (h *HNSW) Save(location string) {
	pLocation := C.CString(location)
	C.saveHNSW(h.index, pLocation)
	C.free(unsafe.Pointer(pLocation))
}

func (h *HNSW) Free() {
	C.freeHNSW(h.index)
}

func normalizeVector(vector []float32) []float32 {
	var norm float32
	for i := 0; i < len(vector); i++ {
		norm += vector[i] * vector[i]
	}
	norm = 1.0 / (float32(math.Sqrt(float64(norm))) + 1e-15)
	for i := 0; i < len(vector); i++ {
		vector[i] = vector[i] * norm
	}
	return vector
}

func (h *HNSW) AddPoint(vector []float32, label uint32) {
	if h.normalize {
		vector = normalizeVector(vector)
	}
	C.addPoint(h.index, (*C.float)(unsafe.Pointer(&vector[0])), C.ulong(label))
}

func (h *HNSW) SearchKNN(vector []float32, N int) ([]uint32, []float32) {
	Clabel := make([]C.ulong, N, N)
	Cdist := make([]C.float, N, N)
	if h.normalize {
		vector = normalizeVector(vector)
	}
	numResult := int(C.searchKnn(h.index, (*C.float)(unsafe.Pointer(&vector[0])), C.int(N), &Clabel[0], &Cdist[0]))
	labels := make([]uint32, N)
	dists := make([]float32, N)
	for i := 0; i < numResult; i++ {
		labels[i] = uint32(Clabel[i])
		dists[i] = float32(Cdist[i])
	}
	return labels[:numResult], dists[:numResult]
}

func (h *HNSW) SetEf(ef int) {
	C.setEf(h.index, C.int(ef))
}
