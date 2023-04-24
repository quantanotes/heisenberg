package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCollection(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/tests/1/")
	db := NewDB(path)
	defer db.Close()

	_, err := db.NewCollection("a", 3, 1000, "cosine", 20, 20)
	if err != nil {
		panic(err)
	}
}

func TestPutGet(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/tests/1/")
	db := NewDB(path)
	defer db.Close()

	err := db.Put("a", []float32{1, 2, 3}, struct{ Data string }{"hello"}, "a")
	if err != nil {
		panic(err)
	}

	v, err := db.Get("a", "a")
	if err != nil || v == nil {
		panic(err)
	}

	fmt.Println(v.Meta)
}

func TestDB(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/tests/1/")
	db := NewDB(path)
	defer db.Close()

	_, err := db.NewCollection("a", 3, 1000, "cosine", 20, 20)
	if err != nil {
		panic(err)
	}
}
