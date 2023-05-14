package main

import (
	"heisenberg/core"
	"heisenberg/server"
	"os"
	"path/filepath"
)

func main() {
	wd, _ := os.Getwd()
	h := core.NewHeisenberg(filepath.Join(wd, "/data/"))
	defer h.Close()
	server.RunAPI(h)
}
