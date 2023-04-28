package store

import (
	"fmt"
	"heisenberg/log"
)

// Interface to handle sharding via consistent hashing
type shard struct {
	shards []string // Shard addresses
	ring   ring     // Circular data structure for consistent hashing
}

func (s *shard) AddShard(addr string) error {

	// TODO: Implement resharding

	s.shards = append(s.shards, addr)
	s.ring.addNode(addr)
	return nil
}

func (s *shard) RemoveShard(addr string) error {

	// TODO: Implement resharding

	for i, sa := range s.shards {
		if sa == addr {
			s.shards = append(s.shards[:i], s.shards[i+1:]...)
			s.ring.removeNode(addr)
			return nil
		}
	}

	return fmt.Errorf("")
}

func (s *shard) GetShard(key []byte) (string, error) {
	addr := s.ring.getNode(key)
	if addr == "" {
		err := "shards size is none"
		log.Error(err, nil)
		return addr, fmt.Errorf(err)
	}
	return addr, nil
}
