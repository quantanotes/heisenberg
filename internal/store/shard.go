package store

import (
	"fmt"
	"heisenberg/log"
)

// Interface to handle sharding via consistent hashing
type shard struct {
	shards   []string             // Shard id
	replicas *map[string]*replica // Shard clients with replication management
	ring     ring                 // Circular data structure for consistent hashing
}

func (s *shard) addShard(id string) error {
	if s.hasShard(id) {
		return fmt.Errorf("")
	}
	old := *s // Create copy of old shard state for resharding
	s.shards = append(s.shards, id)
	s.ring.addNode(id)
	s.reshard(old)
	return nil
}

func (s *shard) removeShard(id string) error {
	old := *s // Create copy of old shard state for resharding
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

func (s *shard) getShard(key []byte) (*replica, error) {
	id := s.ring.getNode(key)
	if id == "" {
		err := "shards size is none"
		log.Error(err, nil)
		return nil, fmt.Errorf(err)
	}
	replica := (*s.replicas)[id]
	return replica, nil
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

func (s *shard) GetReplicas() *map[string]*replica {
	return s.replicas
}
