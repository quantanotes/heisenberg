package store

import "hash/fnv"

func hash(key []byte) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}
