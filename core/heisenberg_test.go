package core

import (
	"fmt"
	"github.com/quantanotes/heisenberg/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestHeisenberg(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/.tmp")
	h := NewHeisenberg(path)
	defer h.Close()
	err := h.NewCollection("c", 3, utils.Cosine)
	if err != nil {
		panic(err)
	}
	err = h.Put("c", "k", []float32{1, 2, 3}, map[string]interface{}{"msg": "bruh"})
	if err != nil {
		panic(err)
	}
	v, err := h.Get("c", "k")
	if err != nil {
		panic(err)
	}
	fmt.Println(v.Value.Meta)
}
