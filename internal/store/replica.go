package store

import (
	"fmt"
	"math/rand"
)

// Interface for handling replication
type replica struct {
	shardId string
	clients *map[string]*StoreClient
}

// Choose random replicas to distribute read requests evenly
func (r *replica) choose() (*StoreClient, error) {
	size := len(*r.clients)
	if size == 0 {
		return nil, fmt.Errorf("no replicas")
	}
	idx := rand.Intn(len(*r.clients))
	var replica *StoreClient
	for _, r := range *r.clients {
		if idx == 0 {
			replica = r
			break
		}
		idx--
	}
	return replica, nil
}
