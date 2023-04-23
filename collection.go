package main

import (
	"fmt"
	"path/filepath"
	"strconv"
)

const configKey = "__HEISENBERG_CONFIG"
const mappingKey = "__HEISENBERG_MAPPING_"

type collectionConfig struct {
	size     int    // max amount of vectors in collectioni
	dim      int    // dimension size of vectors
	space    string // distance measure
	cursor   int    // position to insert vector
	freeList []int  // queue data structure to store deleted vectors
}

type Collection struct {
	name   string // name of bucket in kv store
	mapping string // bucket to map keys to index 
	kv     *KV
	idx    *Index
	config collectionConfig
}

func NewCollection(name string, kv *KV, idx *Index, size int, dim int, space string) (*Collection, error) {
	if err := kv.NewCollection(name); err != nil {
		return nil, err
	}

	config := collectionConfig{
		size,
		dim,
		space,
		0,
		make([]int, 0),
	}

	mapping := fmt.Sprintf("%s%s", mappingKey, name)
	if err := kv.NewCollection(mapping); err != nil {
		return nil, err
	}

	c := &Collection{
		name,
		mapping,	
		kv,
		idx,
		config,
	}

	return c, nil
}

func UseCollection(name string, kv *KV) (*Collection, error) {
	pair, err := kv.Get(configKey, name)
	if err != nil {
		return nil, err
	}

	config := pair.V.Meta.(collectionConfig)

	mapping := fmt.Sprintf("%s%s", mappingKey, name)
	
	idxPath := filepath.Join(kv.path, fmt.Sprintf("%s.idx"), name)
	idx := LoadIndex(idxPath, config.dim, config.space)

	return &Collection{
		name,
		mapping,
		kv,
		idx,
		config,
	}, nil
}

func (c *Collection) DeleteCollection() error {
	err := c.kv.DeleteCollection(c.name)
	err = c.kv.DeleteCollection(c.mapping)
	c = nil
	return err
}

func (c *Collection) Put(p pair) error {
	idx := c.config.cursor
	// set hnsw index to value from free list and pop it
	if len(c.config.freeList) > 0 {
		idx := c.config.freeList[0]
		c.config.freeList := c.config.freeList[1:]
	} else {
		c.config.cursor++
	}

	vec := p.V.Vec

	err := c.kv.Put(p, c.name)
	if err != nil {
		return err
	}
	
	

	err = c.kv.Put()
	c.idx.Put(vec, strconv.)
}

func (c *Collection) Get() {

}

func (c *Collection) Search() {

}

func (c *Collection) Delete() {

}
