package internal

import (
	"hash/fnv"
	"os"
)

var Wd string

func Hash(key []byte) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func Contains[T comparable](arr []T, val T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func init() {
	Wd, _ = os.Getwd()
}
