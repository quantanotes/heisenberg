package store

import (
	"context"
	"fmt"
	"testing"
)

func TestMasterPong(t *testing.T) {
	ctx := context.Background()

	m, _ := NewStoreMasterServer()
	go m.Run(ctx, "localhost:691")
	defer m.Close()

	client, err := NewStoreClient(ctx, "localhost:691")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	pong, err := client.Ping(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)
}

func TestMasterConnect(t *testing.T) {
	ctx := context.Background()

	m, _ := NewStoreMasterServer()
	go m.Run(ctx, "localhost:691")
	defer m.Close()

	mc, _ := NewStoreClient(ctx, "localhost:691")
	defer mc.Close()

	s, _ := NewStoreServer("a")
	go s.Run(ctx, "localhost:692")
	defer s.Close()

	sc, _ := NewStoreClient(ctx, "localhost:692")
	defer sc.Close()
}
