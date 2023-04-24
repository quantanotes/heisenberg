package main

import (
	"os"
	"path/filepath"
)

func main() {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/tests/")
	db := NewDB(path)
	server := NewServer(db)
	server.Run()
}
