package internal

import (
	"hash/fnv"
)

func Hash(key []byte) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}
