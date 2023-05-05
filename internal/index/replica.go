package index

import (
	"fmt"
	"heisenberg/internal"
	"math/rand"
)

// Interface for handling replication
type replica struct {
	Ids     []string
	clients map[string]*IndexClient
}

func (r *replica) addReplica(c *IndexClient, id string) error {
	if c == nil {
		return internal.NilClientError()
	}
	r.clients[id] = c
	if !internal.Contains(r.Ids, id) {
		r.Ids = append(r.Ids, id)
	}
	return nil
}

// Choose random replicas to distribute read requests evenly
func (r *replica) choose() (*IndexClient, error) {
	size := len(r.clients)
	if size == 0 {
		return nil, fmt.Errorf("no replicas")
	}
	idx := rand.Intn(size)
	var replica *IndexClient
	for _, r := range r.clients {
		if idx == 0 {
			replica = r
			break
		}
		idx--
	}
	return replica, nil
}
