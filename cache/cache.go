package cache

import (
	"heisenberg/common"
	"heisenberg/err"
)

type Cache struct {
	buckets map[string]bucket
}

func (c *Cache) NewBucket(name string) error {
	if c.HasBucket(name) {
		return err.BucketExistErr(name)
	}
	c.newBucket(name)
	return nil
}

func (c *Cache) EnsureBucket(name string) {
	if !c.HasBucket(name) {
		bucket := newBucket(name)
		c.buckets[name] = bucket
	}
}

func (c *Cache) DeleteBucket(name string) error {
	if !c.HasBucket(name) {
		return err.BucketNotExistErr(name)
	}
	delete(c.buckets, name)
	return nil
}

func (c *Cache) HasBucket(name string) bool {
	_, ok := c.buckets[name]
	return ok
}

func (c *Cache) newBucket(name string) {
	bucket := newBucket(name)
	c.buckets[name] = bucket
}

func (c *Cache) Get(bucket string, key string) (*common.Value, error) {
	if b, ok := c.buckets[bucket]; ok {
		return b.get(key)
	}
	return nil, err.BucketNotExistErr(bucket)
}

func (c *Cache) Put(bucket string, key string, vector []float32, meta common.Meta) error {
	if b, ok := c.buckets[key]; ok {
		b.put(key, vector, meta)
	}
	return err.BucketNotExistErr(bucket)
}

func (c *Cache) Delete(bucket string, key string) error {
	if b, ok := c.buckets[key]; ok {
		b.delete(key)
	}
	return err.BucketNotExistErr(bucket)
}

func (c *Cache) VectorSearch(bucket string) {}

func (c *Cache) GraphSearch(bucket string) {}

func (c *Cache) DocumentSearch(bucket string) {}
