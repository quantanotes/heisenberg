package master

import (
	"fmt"
	"heisenberg/internal/store"
	"math/rand"
)

// Interface for handling replication
type replica struct {
	clients map[string]*store.StoreClient
}

func (r *replica) addReplica(c *store.StoreClient, id string) {
	r.clients[id] = c
}

// Choose random replicas to distribute read requests evenly
func (r *replica) choose() (*store.StoreClient, error) {
	size := len(r.clients)
	if size == 0 {
		return nil, fmt.Errorf("no replicas")
	}
	idx := rand.Intn(size)
	var replica *store.StoreClient
	for _, r := range r.clients {
		if idx == 0 {
			replica = r
			break
		}
		idx--
	}
	return replica, nil
}
