package query

import (
	"context"
	"heisenberg/internal"
	"heisenberg/internal/pb"
	"heisenberg/log"
	"net"

	"storj.io/drpc/drpcserver"
)

type QueryServer struct {
	server *drpcserver.Server
	lis    *net.Listener
}

func NewQueryServer() (QueryServer, error) {
	q := QueryServer{}
	return q, nil
}

func (q QueryServer) Run(ctx context.Context, addr string) error {
	lis, server, err := internal.NewServer(ctx, addr, q)
	if err != nil {
		log.LogErrReturn("Run", err, q.identity())
	}
	q.server = server
	q.lis = lis
	log.Info("starting query server", q.identity())
	return q.server.Serve(ctx, *q.lis)
}

func (q QueryServer) Close() {
	log.Info("closing query server", q.identity())
}

func (q QueryServer) Ping(ctx context.Context, in *pb.Empty) (*pb.Pong, error) {
	log.Info("recieved ping", q.identity())
	pong := &pb.Pong{
		Id:      "0",
		Master:  true,
		Service: uint32(internal.QueryService),
		Shard:   nil,
	}
	return pong, nil
}

func (q QueryServer) Connect(ctx context.Context, in *pb.Connection) (*pb.Empty, error) {
	return nil, nil
}

func (q QueryServer) AddShard(ctx context.Context, in *pb.Shard) (*pb.Empty, error) {
	return nil, nil
}

func (q QueryServer) RemoveShard(ctx context.Context, in *pb.Shard) (*pb.Empty, error) {
	return nil, nil
}

func (q QueryServer) CreateCollection(ctx context.Context, in *pb.Collection) (*pb.Empty, error) {
	return nil, nil
}

func (q QueryServer) DeleteCollection(ctx context.Context, in *pb.Collection) (*pb.Empty, error) {
	return nil, nil
}

func (q QueryServer) Get(ctx context.Context, in *pb.Key) (*pb.Pair, error) {
	return nil, nil
}

func (q QueryServer) Put(ctx context.Context, in *pb.Item) (*pb.Empty, error) {
	return nil, nil
}

func (q QueryServer) Delete(ctx context.Context, in *pb.Key) (*pb.Empty, error) {
	return nil, nil
}

func (q QueryServer) identity() log.M {
	return internal.Identity(q.lis)
}
