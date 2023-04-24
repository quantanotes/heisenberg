package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDB(t *testing.T) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/tests/1/")
	db := NewDB(path)

	db.NewCollection("a", 3, 1000, "cosine", 20, 20)
}
