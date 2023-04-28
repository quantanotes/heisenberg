package internal

import (
	"log"
	"net"
	"time"

	"storj.io/drpc/drpcconn"
)

// BaseClient is a basic template for all types of clients in the database
type BaseClient struct {
	id		int
	conn   	*net.Conn
	client 	interface{}
}

// NewBaseClient creates a new BaseClient
func NewBaseClient(addr string) (*BaseClient, error) {
	conn, err := ConnectClient(addr)
	if err != nil {
		return nil, err
	}

	return &BaseClient{conn: conn}, nil
}

// Close closes the connection to the database
func (bc *BaseClient) Close() {
	(*bc.conn).Close()
}


