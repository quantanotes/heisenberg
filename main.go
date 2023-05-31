package main

import (
	"os"

	"github.com/quantanotes/heisenberg/core"
	"github.com/quantanotes/heisenberg/server"
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
	db := core.NewDB(path)
	defer db.Close()
	server.RunAPI(db, ":8080")
}
