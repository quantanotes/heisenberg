package hnsw

//#cgo CFLAGS: -I./
//#cgo LDFLAGS: -lhnsw
//#include <stdlib.h>
//#include "hnsw_wrapper.h"
import "C"

type hnsw struct {
	index C.hnsw
}
