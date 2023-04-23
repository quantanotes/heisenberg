package main

// Data instance inside key value store
type data struct {
	Idx  int       `json:"idx"`  // id of vector in search index
	Vec  []float32 `json:"vec"`  // stored vector
	Meta any       `json:"meta"` // metadata associated with vector
}
