package main

import (
	"heisenberg/core"
	"heisenberg/server"
	"os"
)

func main() {
	var path string = os.Getenv("HEISENBERG_PATH")
	if path == "" {
		path = "./.heisenberg/"
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	h := core.NewHeisenberg(path)
	defer h.Close()
	server.RunAPI(h, ":8080")
}
