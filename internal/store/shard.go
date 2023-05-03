package store

import (
	"bytes"
	"encoding/gob"
	"heisenberg/internal"
)

// Interface to handle shards via consistent hashing
type shard struct {
	Replicas map[string]replica // Shard clients with replication management
	Ring     ring               // Circular data structure for consistent hashing
}

func newShard() *shard {
	return &shard{
		Replicas: make(map[string]replica),
		Ring:     newRing(),
	}
}

func (s *shard) getShard(key []byte) (*replica, error) {
	id := s.Ring.getNode(key)
	if id == "" {
		return nil, internal.NoShardsError()
	}
	replica := s.Replicas[id]
	return &replica, nil
}

func (s *shard) addShard(id string) error {
	old := *s
	s.Replicas[id] = replica{
		Ids:     make([]string, 0),
		clients: make(map[string]*StoreClient),
	}
	s.Ring.addNode(id)
	return s.reshard(old)
}

func (s *shard) removeShard(id string) error {
	if !s.hasShard(id) {
		return internal.InvalidShardError(id)
	}
	old := *s
	s.Ring.removeNode(id)
	return s.reshard(old)
}

func (s *shard) hasShard(id string) bool {
	_, ok := s.Replicas[id]
	return ok
}

func (s *shard) reshard(old shard) error {
	return nil
}

func (s *shard) getReplicas() *map[string]replica {
	return &s.Replicas
}

func (s *shard) addReplica(c *StoreClient, id string, shard string) error {
	if c == nil {
		return internal.NilClientError()
	}
	replica, ok := s.Replicas[shard]
	if !ok {
		return internal.InvalidShardError(id)
	}
	replica.addReplica(c, id)
	return nil
}

func (s *shard) toByte() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func shardFromBytes(data []byte) (*shard, error) {
	s := &shard{}
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
