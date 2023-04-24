package storage

import (
	"encoding/binary"
	"fmt"
	"path/filepath"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"go.etcd.io/bbolt"
)

const configKey = "__HEISENBERG_CONFIG"   // In-bucket key for indexing configuration
const mappingKey = "__HEISENBERG_MAPPING" // Nested bucket which maps indices to keys

const collectionMemSize = 3 // Max collection indices in memory

type collectionsMemMap = *orderedmap.OrderedMap[string, *Collection]

type Collection struct {
	name string
	db   *DB
	idx  *Index
}

type DB struct {
	path        string            // Path of all database files
	kv          *bbolt.DB         // On disk data
	collections collectionsMemMap // In memory data for collections/indices
}

func NewDB(path string) *DB {
	kv, err := bbolt.Open(filepath.Join(path, "data.db"), 0666, nil)
	if err != nil {
		panic(err)
	}

	collections := orderedmap.New[string, *Collection]()

	return &DB{
		path,
		kv,
		collections,
	}
}

func (db *DB) Close() {
	// Save and close opened indices
	for pair := db.collections.Oldest(); pair != nil; pair = pair.Next() {
		pair.Value.idx.Save()
		pair.Value.idx.Close()
	}
	db.kv.Close()
}

// Create a new collection of indexed vectors in the database. Returns instance of collection.
func (db *DB) NewCollection(name string, dim int, size int, space string, m int, ef int) (*Collection, error) {
	var collection *Collection

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
		configData, err := ToJson(config)
		if err != nil {
			return err
		}

		// Store config data inside collection bucket
		if err = b.Put([]byte(configKey), configData); err != nil {
			return err
		}

		// Create nested bucket for index to key mappings
		if _, err = b.CreateBucket([]byte(mappingKey)); err != nil {
			return err
		}

		// Create a new index and save at path/collection name.idx
		idxPath := filepath.Join(db.path, fmt.Sprintf("%s.idx", name))
		idx := &Index{
			idxPath,
			nil,
			config,
		}
		idx = NewIndex(idxPath, config, m, ef)

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

	// Add collection to mem map
	db.collections.Set(name, collection)
	db.popCollection()

	return collection, nil
}

// Opens a collection of indexed vectors. Returns instance of collection.
func (db *DB) LoadCollection(name string, loadIdx bool) (*Collection, error) {
	var collection *Collection
	collection = nil

	tx := func(tx *bbolt.Tx) error {
		// Get associated bucket with collection
		b := tx.Bucket([]byte(name))
		if b == nil {
			return fmt.Errorf("collection %s does not exist", name)
		}

		// Get config of collection
		configData := b.Get([]byte(configKey))
		if configData == nil {
			return fmt.Errorf("config for index %s does not exist", name)
		}

		// Read bytes in to config
		config := &indexConfig{}
		if err := FromJson(configData, &config); err != nil || config == nil {
			return fmt.Errorf("invalid config for index %s does not exist", name)
		}

		// Load collection index from path/collection name.idx
		idxPath := filepath.Join(db.path, fmt.Sprintf("%s.idx", name))
		idx := &Index{
			idxPath,
			nil,
			config,
		}
		if loadIdx {
			idx = LoadIndex(idxPath, config)
		}

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
	if loadIdx {
		db.collections.Set(name, collection)
		db.popCollection()
	}

	return collection, nil
}

func (db *DB) GetCollection(name string, loadIdx bool) (*Collection, error) {
	collection, ok := db.collections.Get(name)
	if !ok {
		return db.LoadCollection(name, loadIdx)
	}

	if collection.idx.hnsw == nil && loadIdx {
		idxPath := filepath.Join(db.path, fmt.Sprintf("%s.idx", name))
		collection.idx = LoadIndex(idxPath, collection.idx.config)
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

func (db *DB) DeleteCollection(name string) error {
	return nil
}

func (db *DB) Get(key string, collectionName string) (*Value, error) {
	value := &Value{}

	tx := func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(collectionName))
		if b == nil {
			return fmt.Errorf("collection %s does not exist", collectionName)
		}

		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("key %s does not exist in collection %s", key, collectionName)
		}

		FromJson(data, value)

		return nil
	}

	if err := db.kv.View(tx); err != nil || value == nil {
		return nil, err
	}

	return value, nil
}

func (db *DB) Put(key string, vec []float32, meta interface{}, collectionName string) error {
	tx := func(tx *bbolt.Tx) error {
		collection, b, mapping, err := db.getCollectionBucketMapping(tx, collectionName, false)
		if err != nil {
			return err
		}

		// Value to be inserted at key
		value := &Value{}
		// Previous iteration of data at key
		prev := b.Get([]byte(key))
		if prev != nil {
			FromJson(prev, value)
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
		data, err := ToJson(value)
		if err != nil {
			return err
		}

		err = b.Put([]byte(key), data)
		if err != nil {
			return err
		}

		// Map index to key
		idxInBytes := IntToBytes(value.Idx)
		err = mapping.Put([]byte(key), idxInBytes)
		if err != nil {
			return err
		}

		// Iterate index cursor if all kv mutations are succesful
		collection.idx.MutNext()

		return nil
	}

	return db.kv.Update(tx)
}

func (db *DB) Delete(key string, collectionName string) error {
	tx := func(tx *bbolt.Tx) error {
		collection, b, mapping, err := db.getCollectionBucketMapping(tx, collectionName, false)
		if err != nil {
			return err
		}

		// Delete key from key value store
		if err = b.Delete([]byte(key)); err != nil {
			return err
		}

		// Get previous value at key
		value := &Value{}
		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("key %s does not exist in collection %s", key, collectionName)
		}
		FromJson(data, value)

		// Delete mapping to key
		idxInBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(idxInBytes, uint32(value.Idx))
		if err = mapping.Delete(idxInBytes); err != nil {
			return err
		}

		collection.idx.Delete(value.Idx)

		return nil
	}

	return db.kv.Update(tx)
}

func (db *DB) Search(vec []float32, k int, collectionName string) ([]*Pair, error) {
	var results []*Pair

	tx := func(tx *bbolt.Tx) error {
		collection, b, mapping, err := db.getCollectionBucketMapping(tx, collectionName, true)
		if err != nil {
			return err
		}

		// Retrieve knn results from index
		resultIds := collection.idx.Search(vec, k)

		results := make([]*Pair, len(resultIds))

		// Get vector and meta data for each knn result
		for _, id := range resultIds {
			key := mapping.Get(IntToBytes(int(id)))
			if key == nil {
				break
			}

			value := Value{}
			data := b.Get(key)
			if data == nil {
				break
			}
			err := FromJson(data, value)
			if err != nil {
				break
			}

			results = append(results, &Pair{
				string(key),
				value,
			})
		}

		return nil
	}

	if err := db.kv.View(tx); err != nil {
		return nil, err
	}

	return results, nil
}

// Helpful utility function that pops up in delete and put
func (db *DB) getCollectionBucketMapping(tx *bbolt.Tx, collectionName string, loadIdx bool) (*Collection, *bbolt.Bucket, *bbolt.Bucket, error) {
	collection, err := db.GetCollection(collectionName, loadIdx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("collection %s does not exist", collectionName)
	}

	// Get associated on disk bucket for collection
	b := tx.Bucket([]byte(collectionName))
	if b == nil {
		return nil, nil, nil, fmt.Errorf("collection %s does not exist", collectionName)
	}

	// Retrive index mapping for collection
	mapping := b.Bucket([]byte(mappingKey))
	if mapping == nil {
		return nil, nil, nil, fmt.Errorf("index mapping for collection %s does not exist", collectionName)
	}

	return collection, b, mapping, nil
}
