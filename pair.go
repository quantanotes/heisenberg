package main

// Data instance inside key value store
type Value struct {
	Idx  int         `json:"idx"`  // id of vector in search index
	Vec  []float32   `json:"vec"`  // stored vector
	Meta interface{} `json:"meta"` // metadata associated with vector
}

type Pair struct {
	key   string
	value Value
}
