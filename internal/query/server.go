package internal

import (
	"heisenberg/internal/index"
	"heisenberg/internal/store"
)

type QueryServer struct {
	indexClient *index.IndexClient
	storeClient *store.StoreClient
}

func NewQueryServer(indexAddr string, storeAddr string) *QueryServer {
	indexClient, err := index.NewIndexClient(indexAddr)
	if err != nil {
		return nil
	}

	storeClient, err := store.NewStoreClient(storeAddr)
	if err != nil {
		return nil
	}

	return &QueryServer{
		indexClient,
		storeClient,
	}
}
