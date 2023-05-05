package internal

import (
	"context"
	"net"

	"heisenberg/internal/pb"

	"storj.io/drpc/drpcconn"
)

// Generic interface for clients
type Client struct {
	addr   string
	conn   *drpcconn.Conn
	client pb.DRPCServiceClient
}

func NewClient(ctx context.Context, addr string, service Service) (*Client, error) {
	conn, err := connect(addr)
	if err != nil {
		return nil, err
	}
	client := &Client{
		addr,
		conn,
		pb.NewDRPCServiceClient(conn),
	}
	err = client.validateClient(ctx, service)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func connect(addr string) (*drpcconn.Conn, error) {
	tr, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, ConnectionError(addr, err)
	}
	conn := drpcconn.New(tr)
	return conn, nil
}

// Ensure client is connecting to correct service
func (c *Client) validateClient(ctx context.Context, expected Service) error {
	pong, err := c.client.Ping(ctx, nil)
	if err != nil {
		return ConnectionError(c.addr)
	}
	if pong.Service != uint32(expected) {
		return IncorrectServiceError(expected, Service(pong.Service))
	}
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Connect(ctx context.Context, addr string) {
	c.client.Connect(ctx, &pb.Connection{Address: addr})
}

func (c *Client) AddShard(ctx context.Context, id string) {
	c.client.AddShard(ctx, &pb.Shard{Shard: id})
}

func (c *Client) CreateCollection(ctx context.Context, collection []byte) error {
	_, err := c.client.CreateCollection(ctx, &pb.Collection{Collection: collection})
	return err
}

func (c *Client) DeleteCollection(ctx context.Context, collection []byte) error {
	_, err := c.client.DeleteCollection(ctx, &pb.Collection{Collection: collection})
	return err
}

func (c *Client) Get(ctx context.Context, key []byte, collection []byte) (*pb.Pair, error) {
	pair, err := c.client.Get(ctx, &pb.Key{Key: key, Collection: collection})
	return pair, err
}

func (c *Client) Put(ctx context.Context, key []byte, value []byte, collection []byte) error {
	_, err := c.client.Put(ctx, &pb.Item{Key: key, Value: value, Collection: collection})
	return err
}

func (c *Client) Delete(ctx context.Context, key []byte, collection []byte) error {
	_, err := c.client.Delete(ctx, &pb.Key{Key: key, Collection: collection})
	return err
}

func (c *Client) Ping(ctx context.Context) (*pb.Pong, error) {
	return c.client.Ping(ctx, nil)
}
