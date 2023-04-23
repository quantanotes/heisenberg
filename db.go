package main

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

const configKey = "__HEISENBERG_CONFIG"    // in-bucket key for indexing configuration for collection
const mappingKey = "__HEISENBERG_MAPPING_" // prefix for bucket which maps indexes to keys

type DB struct {
	path string
	kv   *bbolt.DB
}

func NewDB(path string) {

}

func UseDB(path string) {

}

func (db *DB) Close() {
	db.kv.Close()
}

// Create a new collection of indexed vectors in the database. Returns instance of collection.
func (db *DB) NewCollection(name string, dim int, size int, space string) (*Collection, error) {
	var collection *Collection
	collection = nil

	tx := func(tx *bbolt.Tx) error {
		// Create associated bucket with collection
		b, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return err
		}

		// Create config for collection
		config := &collectionConfig{
			size,
			dim,
			space,
			0,
			make([]int, 0),
		}
		configData, err := ToBytes(config)
		if err != nil {
			return err
		}

		// Store config data inside collection bucket
		err = b.Put([]byte(configKey), configData)
		if err != nil {
			return err
		}

		collection = &Collection{
			name,
			db,
			*config,
		}

		return nil
	}

	if err := db.kv.Update(tx); err != nil || collection == nil {
		return nil, err
	} else {
		return collection, nil
	}
}

// Opens a collection of indexed vectors. Returns instance of collection.
func (db *DB) OpenCollection(name string) (*Collection, error) {
	var collection *Collection
	collection = nil

	tx := func(tx *bbolt.Tx) error {
		// Get associated bucket with collection
		b := tx.Bucket([]byte(name))
		if b == nil {
			return errors.New(fmt.Sprintf("collection %s does not exist", name))
		}

		// Get config of collection
		configData := b.Get([]byte(configKey))
		if configData == nil {
			return errors.New(fmt.Sprintf("config for collection %s does not exist", name))
		}

		// Read bytes in to config
		config := collectionConfig{}
		if err := FromBytes(configData, &config); err != nil || config == nil {
			return errors.New(fmt.Sprintf("config for collection %s does not exist", name))
		}

		collection = &Collection{
			name,
			db,
			config,
		}

		return nil
	}

	if err := db.kv.View(tx); err != nil || collection == nil {
		return nil, err
	} else {
		return collection, nil
	}
}
