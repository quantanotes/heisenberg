package internal

import (
	"context"
	"net"

	"storj.io/drpc/drpcconn"
	
	"heisenberg/store"
	"heisenberg/internal"
)

type Client interface {

	Ping(ctx context.Context) Service
}

func NewClient(ctx context.Context, addr string, service Service) (Client, error) {
	conn, err := ConnectClient(addr)
	if err != nil {
		return nil, err
	}

	var client Client
	switch service {
		case QueryService:
			client = nil
		case IndexService:
			client = nil
		case StoreService:
			client = &store.StoreClient{
				conn:   conn,
				client: store.NewDRPCStoreClient(conn),
			}
	default:
		return nil, IncorrectServiceError(service, NoneService)
	}

	err = ValidateClient(client, service, ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ConnectClient(addr string) (*drpcconn.Conn, error) {
	tr, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, ConnectionError(addr, err)
	}
	conn := drpcconn.New(tr)
	return conn, nil
}

func ValidateClient(c Client, expected Service, ctx context.Context) error {
	pong := c.Ping(ctx)
	if pong.Code != expected.Code {
		return IncorrectServiceError(expected, pong)
	}
	return nil
}
