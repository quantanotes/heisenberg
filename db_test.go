package main

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

	_, err := db.NewCollection("a", 3, 1000, "cosine", 20, 20)
	if err != nil {
		panic(err)
	}
}

func TestPutGet(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/tests/1/")
	db := NewDB(path)

	meta := struct{ data string }{data: "hello"}

	err := db.Put("a", []float32{1, 2, 3}, meta, "a")
	if err != nil {
		panic(err)
	}
	v, err := db.Get("a", "a")
	if err != nil {
		panic(err)
	}

	fmt.Println(v.Meta)
}

func TestDB(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/tests/1/")
	db := NewDB(path)

	_, err := db.NewCollection("a", 3, 1000, "cosine", 20, 20)
	if err != nil {
		panic(err)
	}
}
