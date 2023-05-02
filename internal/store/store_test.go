package store

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

var wg sync.WaitGroup // wait for go routines to finish before main thread

func TestMasterPong(t *testing.T) {
	defer wg.Done()
	ctx := context.Background()
	wg.Add(1)
	// create master server
	go RunStoreMasterServer(ctx, "localhost:691")

	client, err := NewStoreClient(ctx, "localhost:691")
	if err != nil {
		panic(err)
	}

	pong, err := client.Ping(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(pong)

	wg.Wait()
}
