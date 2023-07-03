package cache

import (
	"context"
	"heisenberg/common"
	"heisenberg/pb"

	"github.com/hashicorp/raft"
)

type Server struct {
	id    string
	addr  string
	mesh  common.Mesh
	sRaft *raft.Raft // shard raft group
	bRaft *raft.Raft // raft for bucket consensus (global to all cache nodes)
}

func NewServer() *Server {
	return nil
}

func (s *Server) Serve(ctx context.Context) {
	go common.Serve[pb.DRPCCacheServiceServer](ctx, s.addr, pb.DRPCRegisterCacheService, s)
}

func (s *Server) Get(ctx context.Context, key *pb.Key) (*pb.Value, error) {
	// TODO: route to master shard
	return nil, nil
}

func (s *Server) Put(ctx context.Context, entry *pb.Entry) (*pb.Nil, error) {
	_, master := s.sRaft.LeaderWithID()
	if s.id != string(master) {
		return nil, nil
	}
	payload := shardFSMPayload{
		Cmd:    "PUT",
		Bucket: entry.Bucket,
		Key:    entry.Key,
		Vector: entry.Vector,
		Meta:   entry.Meta.AsMap(),
	}
	common.ApplyRaft(s.sRaft, payload)
	return nil, nil
}

func (s *Server) Delete(ctx context.Context, key *pb.Key) (*pb.Nil, error) {
	_, master := s.sRaft.LeaderWithID()
	if s.id != string(master) {
		return nil, nil
	}
	payload := shardFSMPayload{
		Cmd:    "DELETE",
		Bucket: key.Bucket,
		Key:    key.Key,
	}
	common.ApplyRaft(s.sRaft, payload)
	return nil, nil
}
