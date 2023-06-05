package core

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/quantanotes/heisenberg/common"

	"go.etcd.io/bbolt"
)

const (
	collectionPrefix = "__" // All collection names are prefixed to distinguish from mapping buckets.
	configKey        = "config"
)

type DB struct {
	kv *bbolt.DB     // key value store to hold meta data
	im *IndexManager // manages memory for indices
}

func NewDB(path string) *DB {
	kv, err := bbolt.Open(filepath.Join(path, "data.db"), 0600, nil)
	if err != nil {
		panic(err)
	}
	im := NewIndexManager(path, 100)
	return &DB{kv, im}
}

func (db *DB) Close() {
	db.im.Close()
	db.kv.Close()
}

func (db *DB) NewCollection(collection string, dim uint, space common.SpaceType) error {
	tx := func(tx *bbolt.Tx) error {
		// Create collection bucket
		_, err := tx.CreateBucketIfNotExists([]byte(collectionPrefix + collection))
		if err != nil {
			return err
		}
		// Create mapping bucket
		m, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		// Store index configuration in mapping bucket
		conf := common.IndexConfig{Name: collection, Indexer: int(HNSWIndexerType), FreeList: make([]uint, 0), Dim: dim, Space: int(space), Count: 0}
		b, err := common.ToBytes(conf)
		if err != nil {
			return err
		}
		// Create new index
		m.Put([]byte(configKey), b)
		db.im.New(collection, HNSWIndexerType, dim, space)
		return nil
	}
	return db.kv.Update(tx)
}

func (db *DB) DeleteCollection(collection string) error {
	tx := func(tx *bbolt.Tx) error {
		// Delete collection bucket
		err := tx.DeleteBucket([]byte(collectionPrefix + collection))
		if err != nil {
			return err
		}
		// Delete mapping bucket
		err = tx.DeleteBucket([]byte(collection))
		if err != nil {
			return err
		}
		// Delete index
		err = db.im.Delete(collection)
		if err != nil {
			return err
		}
		return err
	}
	return db.kv.Update(tx)
}

func (db *DB) Get(collection string, key string) (Entry, error) {
	cb := []byte(collectionPrefix + collection)
	kb := []byte(key)
	var data []byte
	tx := func(tx *bbolt.Tx) error {
		// Retrieve bucket
		b := tx.Bucket(cb)
		if b == nil {
			return common.InvalidCollection(collection)
		}
		// Retrieve value at key
		data = b.Get(kb)
		if data == nil {
			return common.InvalidKey(key, collection)
		}
		return nil
	}
	// Execute transaction
	if err := db.kv.View(tx); err != nil {
		return Entry{}, err
	}
	// Deserialise raw value
	val, err := DeserialiseValue(data)
	if err != nil {
		return Entry{}, err
	}
	return Entry{collection, key, val}, nil
}

func (db *DB) Put(collection string, key string, vec []float32, meta map[string]any) error {
	kb := []byte(key)
	tx := func(tx *bbolt.Tx) error {
		// Retrieve bucket, mapping and index of collection
		b, m, i, err := db.getBucketMappingIndex(tx, collection)
		if err != nil {
			return err
		}
		val := Value{}
		prev := b.Get(kb) // Previous value at key
		if prev != nil {
			val, err = DeserialiseValue(prev)
			if err != nil {
				return err
			}
			// Update value
			if vec != nil {
				val.Vector = vec
			}
			if meta != nil {
				val.Meta = meta
			}
		} else {
			if vec == nil {
				return errors.New("vector cannot be empty")
			}
			val.Index = i.Next()
			val.Vector = vec
			val.Meta = meta
			// Store index-key mapping in mapping bucket
			m.Put(common.IntToBytes(int(val.Index)), kb)
		}
		data, err := val.Serialise()
		if err != nil {
			return err
		}
		// Store value at key
		b.Put(kb, data)
		// Insert vector to index
		if vec != nil {
			err = i.Insert(val.Index, vec)
			if err != nil {
				return err
			}
		}
		go i.Save(db.im.GetPath(collection))
		return nil
	}
	return db.kv.Update(tx)
}

func (db *DB) Delete(collection string, key string) error {
	kb := []byte(key)
	tx := func(tx *bbolt.Tx) error {
		// Retrieve bucket, mapping and index of collection
		b, m, i, err := db.getBucketMappingIndex(tx, collection)
		if err != nil {
			return err
		}
		data := b.Get(kb) // Previous value	at key
		if data == nil {
			return common.InvalidKey(key, collection)
		}
		val, err := DeserialiseValue(data)
		if err != nil {
			return err
		}
		// Delete value at key
		if err = b.Delete([]byte(key)); err != nil {
			return err
		}
		// Delete key-index mapping from mapping bucket
		if err = m.Delete(common.IntToBytes(int(val.Index))); err != nil {
			return err
		}
		// Delete vector from index
		if err := i.Delete(val.Index); err != nil {
			return err
		}
		go i.Save(db.im.GetPath(collection))
		return nil
	}
	return db.kv.Update(tx)
}

func (db *DB) Search(collection string, query []float32, k uint) ([]Entry, error) {
	var results []Entry
	tx := func(tx *bbolt.Tx) error {
		// Retrieve bucket, mapping and index of collection
		b, m, i, err := db.getBucketMappingIndex(tx, collection)
		if err != nil {
			return err
		}
		// Perform KNN search and retrieve indices
		ids, err := i.Search(query, k)
		if err != nil {
			return err
		}
		// Retrieve values
		results = make([]Entry, 0)
		for _, id := range ids {
			// Retrieve mapping to value
			key := m.Get(common.IntToBytes(int(id)))
			if key == nil {
				continue
			}
			// Retrive value
			data := b.Get(key)
			if data == nil {
				continue
			}
			val, err := DeserialiseValue(data)
			if err != nil {
				continue
			}
			results = append(results, Entry{
				collection,
				string(key),
				val,
			})
		}
		return nil
	}
	// Execute transaction
	if err := db.kv.View(tx); err != nil {
		return nil, err
	}
	return results, nil
}

// Retrieves bucket, mapping and index for a given collection.
func (db *DB) getBucketMappingIndex(tx *bbolt.Tx, collection string) (*bbolt.Bucket, *bbolt.Bucket, Index, error) {
	// Retrieve collection bucket
	b := tx.Bucket([]byte(collectionPrefix + collection))
	if b == nil {
		return nil, nil, nil, common.InvalidCollection(collection)
	}
	// Retrieve mapping from keys to index for collection
	m := tx.Bucket([]byte(collection))
	if m == nil {
		return nil, nil, nil, fmt.Errorf("key-index mapping for collection %s does not exist", collection)
	}
	// Retrieve index
	idx, err := db.im.Get(collection, db.kv)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("index for collection %s does not exist, trace: %s", collection, err.Error())
	}
	return b, m, idx, nil
}
