package store

import (
	"fmt"
	"heisenberg/internal"

	"storj.io/drpc/drpcconn"
)

type StoreClient struct {
	conn *drpcconn.Conn
}

func NewStoreClient(addr string) (*StoreClient, error) {
	conn, err := internal.NewConnection(addr)
	if err != nil {
		return nil, fmt.Errorf("NewStoreClient error: %v", err)
	}

	return &StoreClient{conn}, nil
}

func (sc *StoreClient) Close() {
	sc.conn.Close()
}

func (sc *StoreClient) Ping() error {
	sc.conn.Close()
	return nil
}
