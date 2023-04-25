package main

import (
	"heisenberg/server"
	"heisenberg/storage"
	"os"
	"path/filepath"
)

func main() {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "/datafiles/")
	db := storage.NewDB(path)
	defer db.Close()
	server := server.NewServer(db)
	server.Run()
}
