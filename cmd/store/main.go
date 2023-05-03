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
	dirPtr := flag.String("dir", "", "directory of storage data files")
	idPtr := flag.String("id", "", "id of node")
	ctx := context.Background()

	flag.Parse()

	log.Info("@store/Main", nil)

	if *hostPtr == "" {
		log.Fatal("host not specified", nil)
	}

	if *dirPtr == "" {
		log.Fatal("directory not specified", nil)
	}

	if *masterPtr {
		m, err := store.NewStoreMasterServer(*dirPtr)
		if err != nil {
			log.Fatal(err.Error(), nil)
			panic(nil)
		}
		defer m.Close()
		m.Run(ctx, *hostPtr)
	} else {
		s, err := store.NewStoreServer(*dirPtr, *idPtr)
		if err != nil {
			log.Fatal(err.Error(), nil)
			panic(nil)
		}
		defer s.Close()
		s.Run(ctx, *hostPtr)
	}
}
