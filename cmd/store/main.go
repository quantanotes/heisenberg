package main

import (
	"context"
	"flag"
	"heisenberg/internal/store"
	"heisenberg/internal/store/master"
	"heisenberg/log"
)

func main() {
	masterPtr := flag.Bool("m", false, "set if master node")
	hostPtr := flag.String("h", "", "host name for node")
	ctx := context.Background()

	flag.Parse()

	log.Info("@Main Heisenberg Store", nil)

	if *hostPtr == "" {
		log.Fatal("Host not specified", nil)
	}

	if *masterPtr {
		master.RunStoreMasterServer(ctx, *hostPtr)
	} else {
		store.RunStoreServer(ctx, *hostPtr)
	}
}
