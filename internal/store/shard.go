package store

import (
	"fmt"
	"heisenberg/log"
)

// Interface to handle sharding via consistent hashing
type shard struct {
	shards  []string                 // Shard id
	ring    ring                     // Circular data structure for consistent hashing
	clients *map[string]*StoreClient // Shard clients
}

func (s *shard) addShard(id string) error {
	if s.hasShard(id) {
		return fmt.Errorf("")
	}
	old := *s
	s.shards = append(s.shards, id)
	s.ring.addNode(id)
	s.reshard(old)
	return nil
}

func (s *shard) removeShard(id string) error {
	old := *s
	for i, sid := range s.shards {
		if sid == id {
			s.shards = append(s.shards[:i], s.shards[i+1:]...)
			s.ring.removeNode(id)
			s.reshard(old)
			return nil
		}
	}
	return fmt.Errorf("")
}

func (s *shard) getShard(key []byte) (string, error) {
	id := s.ring.getNode(key)
	if id == "" {
		err := "shards size is none"
		log.Error(err, nil)
		return id, fmt.Errorf(err)
	}
	return id, nil
}

func (s *shard) hasShard(id string) bool {
	for _, sid := range s.shards {
		if sid == id {
			return true
		}
	}
	return false
}

func (s *shard) reshard(old shard) error {
	return nil
}
