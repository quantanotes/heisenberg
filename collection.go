package main

type collectionConfig struct {
	Size     int    `json:"size"`     // max amount of vectors in collectioni
	Dim      int    `json:"dim"`      // dimension size of vectors
	Space    string `json:"space"`    // distance measure
	Cursor   int    `json:"cursor"`   // position to insert vector
	FreeList []int  `json:"FreeList"` // queue data structure to store deleted vectors
}

type Collection struct {
	name   string
	db     *DB
	config collectionConfig
}
