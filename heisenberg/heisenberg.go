package main

// Master service for distributed database
type struct Heisenberg {
	shards []string // Shard locations of master replica
	replicas []string // Replica locations
}