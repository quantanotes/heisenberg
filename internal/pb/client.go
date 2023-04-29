package internal

import (
	"net"
	"context"

	"storj.io/drpc/drpcconn"
)

// BaseClient is a basic template for all types of clients in the database
type BaseClient struct {
	conn   		*drpcconn.Conn
	client 		DRPCClient
}

// NewBaseClient creates a new BaseClient
func NewBaseClient(ctx context.Context, addr string, service Service, createDrpcClient func(conn *drpcconn.Conn) interface{}) (*BaseClient, error) {
	// Defien the new baseclient return 
	bc := &BaseClient{}

	// Get the connection
	conn, err := ConnectClient(addr)
	if err != nil {
		return nil, err
	}
	bc.conn = conn

	// Define the client
	bc.client = createDrpcClient(conn)

	// Ping the server
	err = ValidateClient(bc, service, ctx)
	if err != nil {
		return nil, err
	}
	return bc, nil	
}

// Close closes the connection to the database
func (bc *BaseClient) Close() {
	(*bc.conn).Close()
}

// Wrapper for the Ping function
func (bc *BaseClient) Ping(ctx context.Context) Service {
	pong, err := bc.client.Ping(ctx, nil)
	if err != nil {
		return NoneService
	}
	return Service{Code: pong.Service.Code, Name: pong.Service.Name}
}

// Dial up connection to the database
func ConnectClient(addr string) (*drpcconn.Conn, error) {
	tr, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, ConnectionError(addr, err)
	}
	conn := drpcconn.New(tr)
	return conn, nil
}

// ValidateClient checks if the client is valid by pinging
func ValidateClient(c Client, expected Service, ctx context.Context) error {
	pong := c.Ping(ctx)
	if pong.Code != expected.Code {
		return IncorrectServiceError(expected, pong)
	}
	return nil
}
