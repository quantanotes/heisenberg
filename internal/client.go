package internal

import (
	"context"
	"net"

	"heisenberg/internal/pb"

	"storj.io/drpc/drpcconn"
)

type Client struct {
	conn   *drpcconn.Conn
	client pb.DRPCServiceClient
}

// Base constructor for clients
func NewClient(ctx context.Context, addr string, service Service) (*Client, error) {
	conn, err := connect(addr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn,
		pb.NewDRPCServiceClient(conn),
	}

	err = client.validateClient(ctx, service)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func connect(addr string) (*drpcconn.Conn, error) {
	tr, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, ConnectionError(addr, err)
	}
	conn := drpcconn.New(tr)
	return conn, nil
}

func (c *Client) validateClient(ctx context.Context, expected Service) error {
	pong := c.ping(ctx)
	if pong.Code != expected.Code {
		return IncorrectServiceError(expected, pong)
	}
	return nil
}

func (c *Client) ping(ctx context.Context) Service {
	pong, err := c.client.Ping(ctx, nil)
	if err != nil {
		return NoneService
	}
	return Service{Code: pong.Service.Code, Name: pong.Service.Name}
}
