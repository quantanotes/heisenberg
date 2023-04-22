package main

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type DB struct {
	kv    *bolt.DB // key-value storage
	index Index    // spatial indexing
}

func NewDB(path string) (*DB, error) {
	kv, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	index := Index{}

	return &DB{
		kv:    kv,
		index: index,
	}, nil
}

func (db *DB) Close() {
	db.kv.Close()
}

func (db *DB) NewCollection(name string) error {
	tx := func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		return err
	}
	return db.kv.Update(tx)
}

func (db *DB) DeleteCollection(name string) error {
	tx := func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(name))
		return err
	}
	return db.kv.Update(tx)
}

func (db *DB) Put(p pair, collection string) error {
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

	err := db.kv.Update(tx)

	return err
}

func (db *DB) Get(key string, collection string) (*pair, error) {
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

	err := db.kv.View(tx)

	return p, err
}

func (db *DB) Delete(key string, collection string) error {
	tx := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if b == nil {
			return fmt.Errorf("collection %v does not exist", collection)
		}

		err := b.Delete([]byte(key))

		return err
	}

	err := db.kv.Update(tx)

	return err
}
