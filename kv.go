package main

import (
	"fmt"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

// Key-value storage
type KV struct {
	path string
	kv   *bolt.DB
}

func NewKV(path string) (*KV, error) {
	kv := &KV{
		path: path,
	}

	bolt, err := bolt.Open(filepath.Join(path, "data.bin"), 0666, nil)
	if err != nil {
		return nil, err
	}

	kv.kv = bolt
	return kv, nil
}

func (kv *KV) Close() {
	kv.kv.Close()
}

func (kv *KV) NewCollection(name string) error {
	tx := func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		return err
	}
	return kv.kv.Update(tx)
}

func (kv *KV) DeleteCollection(name string) error {
	tx := func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(name))
		return err
	}
	return kv.kv.Update(tx)
}

func (kv *KV) Put(p pair, collection string) error {
	tx := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if b == nil {
			return fmt.Errorf("collection %v does not exist", collection)
		}

		k := p.K
		v, err := p.V.toBytes()

		if err != nil {
			return err
		}

		err = b.Put([]byte(k), v)
		if err != nil {
			return err
		}

		return nil
	}

	err := kv.kv.Update(tx)

	return err
}

func (kv *KV) Get(key string, collection string) (*pair, error) {
	p := &pair{K: key}

	tx := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if b == nil {
			return fmt.Errorf("collection %v does not exist", collection)
		}

		d := b.Get([]byte(key))
		if d == nil {
			return fmt.Errorf("key %v does not exist in collection %v", key, collection)
		}

		v, err := valueFromBytes(d)
		if err != nil {
			return err
		}

		p.V = *v

		return nil
	}

	err := kv.kv.View(tx)

	return p, err
}

func (kv *KV) Delete(key string, collection string) error {
	tx := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if b == nil {
			return fmt.Errorf("collection %v does not exist", collection)
		}

		err := b.Delete([]byte(key))

		return err
	}

	err := kv.kv.Update(tx)

	return err
}
