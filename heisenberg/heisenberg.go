package main

// Master service for distributed database
type Heisenberg struct {
	shards    []string // Shard locations of master replica
	replicas  []string // Replica locations
	queryPool []byte   // Don't go in there ;)
}
