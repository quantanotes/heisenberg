package main

// Master service for distributed database
type Heisenberg struct {
	replicas  []string // Replica locations
	queryPool []byte   // Don't go in there ;)
}
