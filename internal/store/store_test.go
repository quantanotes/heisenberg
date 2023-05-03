package store

import (
	"context"
	"fmt"
	"heisenberg/internal"
	"path/filepath"
	"testing"
)

func TestMasterPong(t *testing.T) {
	ctx := context.Background()

	dir := filepath.Join(internal.Wd, "tests/master.bin")
	m, err := NewStoreMasterServer(dir)
	if err != nil {
		panic(err)
	}
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

	dir := filepath.Join(internal.Wd, "tests/master.bin")

	m, _ := LoadStoreMasterServer(dir)
	go m.Run(ctx, "localhost:691")
	defer m.Close()

	mc, _ := NewStoreClient(ctx, "localhost:691")
	defer mc.Close()

	mc.AddShard(ctx, "a")

	dir = filepath.Join(internal.Wd, "tests/replica_a.bin")
	s, _ := NewStoreServer(dir, "a")
	go s.Run(ctx, "localhost:692")
	defer s.Close()

	sc, _ := NewStoreClient(ctx, "localhost:692")
	defer sc.Close()

	mc.Connect(ctx, "localhost:692")
}

func TestStorePutGet(t *testing.T) {
	ctx := context.Background()

	dir := filepath.Join(internal.Wd, "tests/master.bin")
	m, _ := LoadStoreMasterServer(dir)
	go m.Run(ctx, "localhost:691")
	defer m.Close()

	mc, _ := NewStoreClient(ctx, "localhost:691")
	defer mc.Close()

	mc.AddShard(ctx, "a")

	dir = filepath.Join(internal.Wd, "tests/replica_a.bin")
	s, _ := NewStoreServer(dir, "a")
	go s.Run(ctx, "localhost:692")
	defer s.Close()

	sc, _ := NewStoreClient(ctx, "localhost:692")
	defer sc.Close()

	mc.Connect(ctx, "localhost:692")
	mc.CreateCollection(ctx, []byte("collection"))
	mc.Put(ctx, []byte("key"), []byte("value"), []byte("collection"))
	value, _ := mc.Get(ctx, []byte("key"), []byte("collection"))
	fmt.Println(value)
}
