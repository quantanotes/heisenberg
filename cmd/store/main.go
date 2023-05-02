package main

import (
	"context"
	"flag"
	"heisenberg/internal/store"
	"heisenberg/log"
)

func main() {
	masterPtr := flag.Bool("m", false, "set if master node")
	hostPtr := flag.String("h", "", "host name for node")
	//idPtr := flag.String("id", "", "id of node")
	ctx := context.Background()

	flag.Parse()

	log.Info("@Main Heisenberg Store", nil)

	if *hostPtr == "" {
		log.Fatal("Host not specified", nil)
	}

	if *masterPtr {
		m, err := store.NewStoreMasterServer()
		if err != nil {
			log.Fatal(err.Error(), nil)
			panic(nil)
		}
		defer m.Close()
		m.Run(ctx, *hostPtr)
	} else {
		s, err := store.NewStoreServer("a")
		if err != nil {
			log.Fatal(err.Error(), nil)
			panic(nil)
		}
		defer s.Close()
		s.Run(ctx, *hostPtr)
	}
}
