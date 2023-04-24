package main

import (
	"heisenberg/server"
	"heisenberg/storage"
	"os"
	"path/filepath"
)

func main() {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/tests/")
	db := storage.NewDB(path)
	server := server.NewServer(db)
	server.Run()
}
