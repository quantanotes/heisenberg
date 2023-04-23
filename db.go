package main

import (
	"errors"
	"fmt"
	"path/filepath"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"go.etcd.io/bbolt"
)

const configKey = "__HEISENBERG_CONFIG"   // In-bucket key for indexing configuration
const mappingKey = "__HEISENBERG_MAPPING" // Nested bucket which maps indices to keys

const collectionMemSize = 3 // Max collection indices in memory

type collectionsMemMap = orderedmap.OrderedMap[string, *Collection]

type DB struct {
	path        string            // Path of all database files
	kv          *bbolt.DB         // On disk data
	collections collectionsMemMap // In memory data for collections/indices
}

func NewDB(path string) {

}

func LoadDB(path string) {

}

func (db *DB) Close() {
	db.kv.Close()
}

// Create a new collection of indexed vectors in the database. Returns instance of collection.
func (db *DB) NewCollection(name string, dim int, size int, space string, m int, ef int) (*Collection, error) {
	var collection *Collection
	collection = nil

	tx := func(tx *bbolt.Tx) error {
		// Create associated bucket with collection
		b, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return err
		}

		// Create config for collection
		config := &indexConfig{
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

		// Create nested bucket for index to key mappings
		_, err = b.CreateBucket([]byte(mappingKey))
		if err != nil {
			return err
		}

		// Create a new index and save at path/collection name.idx
		idx := NewIndex(filepath.Join(db.path, fmt.Sprintf("%s.idx", name)), config, m, ef)

		collection = &Collection{
			name,
			db,
			idx,
		}

		return nil
	}

	if err := db.kv.Update(tx); err != nil || collection == nil {
		return nil, err
	}

	db.collections.Set(name, collection) // Add collection to mem map
	db.popCollection()

	return collection, nil
}

// Opens a collection of indexed vectors. Returns instance of collection.
func (db *DB) LoadCollection(name string) (*Collection, error) {
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
			return errors.New(fmt.Sprintf("config for index %s does not exist", name))
		}

		// Read bytes in to config
		config := &indexConfig{}
		if err := FromBytes(configData, &config); err != nil || config == nil {
			return errors.New(fmt.Sprintf("invalid config for index %s does not exist", name))
		}

		// Load collection from path/collection name.idx
		idx := LoadIndex(filepath.Join(db.path, fmt.Sprintf("%s.idx", name)), config)

		collection = &Collection{
			name,
			db,
			idx,
		}

		return nil
	}

	if err := db.kv.View(tx); err != nil || collection == nil {
		return nil, err
	}

	// Add collection to mem map
	db.collections.Set(name, collection)
	db.popCollection()

	return collection, nil
}

func (db *DB) GetCollection(name string) (*Collection, error) {
	collection, ok := db.collections.Get(name)
	if !ok {
		return db.LoadCollection(name)
	}
	return collection, nil
}

// Delete oldest collection from collection memmap if overflow
func (db *DB) popCollection() {
	if db.collections.Len() > collectionMemSize {
		key := db.collections.Oldest().Key
		db.CloseCollection(key)
	}
}

// Free index memory for given collection
func (db *DB) CloseCollection(name string) {
	collection, ok := db.collections.Get(name)
	if !ok {
		return
	}
	collection.idx.Save()
	collection.idx.Close()
	db.collections.Delete(name)
}

func (db *DB) Get() *data {
	return nil
}

func (db *DB) Put(key string, vec []float32, meta interface{}, collectionName string) error {
	tx := func(tx *bbolt.Tx) error {
		// Retrieve collection and index
		collection, err := db.GetCollection(collectionName)
		if err != nil {
			return errors.New(fmt.Sprintf("collection %s does not exist", collectionName))
		}

		// Get associated on disk bucket for collection
		b := tx.Bucket([]byte(collectionName))
		if b == nil {
			return errors.New(fmt.Sprintf("collection %s does not exist", collectionName))
		}

		// Value to be inserted at key
		value := &data{}
		// Previous iteration of data at key
		prev := b.Get([]byte(key))
		if prev != nil {
			FromBytes(prev, value)
			if vec != nil {
				value.Vec = vec
			}
			if meta != nil {
				value.Meta = meta
			}
		} else {
			value.Idx = collection.idx.Next()
			value.Vec = vec
			value.Meta = meta
		}

		// Convert to-stored value in to bytes
		data, err := ToBytes(value)
		if err != nil {
			return err
		}

		err = b.Put([]byte(key), data)
		if err != nil {
			return err
		}

		// Iterate index cursor if all kv mutations are succesful
		collection.idx.MutNext()

		return nil
	}

	return db.kv.Update(tx)
}
