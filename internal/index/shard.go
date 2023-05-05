package index

import (
	"bytes"
	"encoding/gob"
	"heisenberg/internal"
	"math"
)

// Interface to handle shards via k-means
type shard struct {
	Replicas  map[string]replica                 // Shard clients with replication management
	Centroids map[string][]float32               // Data structure to handle clustering
	dist      func([]float32, []float32) float32 // Metric
}

func newShard() *shard {
	return &shard{
		Replicas:  make(map[string]replica),
		Centroids: make(map[string][]float32),
	}
}

func (s *shard) centroidDistFunc(vec []float32) func(string) float32 {
	return func(id string) float32 {
		centroid, ok := s.Centroids[id]
		if !ok {
			return math.MaxFloat32
		}
		return s.dist(vec, centroid)
	}
}

func (s *shard) getIds() []string {
	var ids []string
	for id := range s.Replicas {
		ids = append(ids, id)
	}
	return ids
}

func (s *shard) getShards(size int, vec []float32) []*replica {
	dq := newDistqueue(int(size), s.centroidDistFunc(vec))
	dq.pushMany(s.getIds())
	items := make([]*replica, size)
	for i, id := range dq.items {
		item, ok := s.Replicas[id]
		if ok {
			items[i] = &item
		}
	}
	return items
}

func (s *shard) addShard(id string) error {
	old := *s
	s.Replicas[id] = replica{
		Ids:     make([]string, 0),
		clients: make(map[string]*IndexClient),
	}
	return s.reshard(old)
}

func (s *shard) removeShard(id string) error {
	if !s.hasShard(id) {
		return internal.InvalidShardError(id)
	}
	old := *s
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

func (s *shard) addReplica(c *IndexClient, id string, shard string) error {
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
