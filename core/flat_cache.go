package core

import (
	"heisenberg/utils"
)

type flatCache struct {
	name  string
	kv    map[string]Value
	space utils.SpaceType
}

func NewFlatCache(name string, dim uint, space utils.SpaceType) *flatCache {
	return &flatCache{
		name:  name,
		kv:    make(map[string]Value),
		space: space,
	}
}

func (c *flatCache) Get(key string) (Entry, error) {
	if val, ok := c.kv[key]; ok {
		return Entry{Collection: c.name, Key: key, Value: val}, nil
	}
	return Entry{}, utils.InvalidKey(c.name, key)
}

func (c *flatCache) Put(key string, index uint, vector []float32, meta map[string]any) error {
	c.kv[key] = Value{index, vector, meta}
	return nil
}

func (c *flatCache) Delete(key string) error {
	_, ok := c.kv[key]
	if !ok {
		return utils.InvalidKey(key, c.name)
	}
	delete(c.kv, key)
	return nil
}

// TODO
func (c *flatCache) Search(query []float32, k uint) ([]Entry, error) {
	// dq := newDistqueue[[]float32](k, func(vec []float32) float32 {
	// 	switch c.space {
	// 	case utils.Cosine:
	// 		return math.Cosine(query, vec)
	// 	case utils.Ip:
	// 		return math.Ip(query, vec)
	// 	case utils.L2:
	// 		return math.L2(query, vec)
	// 	default:
	// 		panic("invalid space for flat cache")
	// 	}
	// })
	// for _, _ := range c.kv {

	// }

	return nil, nil
}
