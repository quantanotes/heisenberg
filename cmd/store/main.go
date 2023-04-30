package store

import (
	"context"
	"flag"
	"heisenberg/internal/store"
	"heisenberg/internal/store/master"
)

func main() {
	masterPtr := flag.Bool("m", false, "set if master node")
	hostPtr := flag.String("host", "", "host name for node")
	ctx := context.Background()

	flag.Parse()

	if *masterPtr {
		master.RunStoreMasterServer(ctx, *hostPtr)
	} else {
		store.RunStoreServer(ctx, *hostPtr)
	}
}
