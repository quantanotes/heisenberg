package core

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/quantanotes/heisenberg/common"

	"go.etcd.io/bbolt"
)

const (
	bucketPrefix = "__" // All bucket names are prefixed to distinguish from mapping buckets.
	configKey    = "config"
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

func (db *DB) NewBucket(bucket string, dim uint, space common.SpaceType) error {
	tx := func(tx *bbolt.Tx) error {
		// Create bucket bucket
		_, err := tx.CreateBucketIfNotExists([]byte(bucketPrefix + bucket))
		if err != nil {
			return err
		}
		// Create mapping bucket
		m, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		// Store index configuration in mapping bucket
		conf := common.IndexConfig{Name: bucket, Indexer: int(HNSWIndexerType), FreeList: make([]uint, 0), Dim: dim, Space: int(space), Count: 0}
		b, err := common.ToBytes(conf)
		if err != nil {
			return err
		}
		// Create new index
		m.Put([]byte(configKey), b)
		db.im.New(bucket, HNSWIndexerType, dim, space)
		return nil
	}
	return db.kv.Update(tx)
}

func (db *DB) DeleteBucket(bucket string) error {
	tx := func(tx *bbolt.Tx) error {
		// Delete bucket bucket
		err := tx.DeleteBucket([]byte(bucketPrefix + bucket))
		if err != nil {
			return err
		}
		// Delete mapping bucket
		err = tx.DeleteBucket([]byte(bucket))
		if err != nil {
			return err
		}
		// Delete index
		err = db.im.Delete(bucket)
		if err != nil {
			return err
		}
		return err
	}
	return db.kv.Update(tx)
}

func (db *DB) Get(bucket string, key string) (Entry, error) {
	cb := []byte(bucketPrefix + bucket)
	kb := []byte(key)
	var data []byte
	tx := func(tx *bbolt.Tx) error {
		// Retrieve bucket
		b := tx.Bucket(cb)
		if b == nil {
			return common.InvalidBucket(bucket)
		}
		// Retrieve value at key
		data = b.Get(kb)
		if data == nil {
			return common.InvalidKey(key, bucket)
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
	return Entry{bucket, key, val}, nil
}

func (db *DB) Put(bucket string, key string, vec []float32, meta map[string]any) error {
	kb := []byte(key)
	tx := func(tx *bbolt.Tx) error {
		// Retrieve bucket, mapping and index of bucket
		b, m, i, err := db.getBucketMappingIndex(tx, bucket)
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
		go i.Save(db.im.GetPath(bucket))
		return nil
	}
	return db.kv.Update(tx)
}

func (db *DB) Delete(bucket string, key string) error {
	kb := []byte(key)
	tx := func(tx *bbolt.Tx) error {
		// Retrieve bucket, mapping and index of bucket
		b, m, i, err := db.getBucketMappingIndex(tx, bucket)
		if err != nil {
			return err
		}
		data := b.Get(kb) // Previous value	at key
		if data == nil {
			return common.InvalidKey(key, bucket)
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
		go i.Save(db.im.GetPath(bucket))
		return nil
	}
	return db.kv.Update(tx)
}

func (db *DB) Search(bucket string, query []float32, k uint) ([]Entry, error) {
	var results []Entry
	tx := func(tx *bbolt.Tx) error {
		// Retrieve bucket, mapping and index of bucket
		b, m, i, err := db.getBucketMappingIndex(tx, bucket)
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
				bucket,
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

// Retrieves bucket, mapping and index for a given bucket.
func (db *DB) getBucketMappingIndex(tx *bbolt.Tx, bucket string) (*bbolt.Bucket, *bbolt.Bucket, Index, error) {
	// Retrieve bucket bucket
	b := tx.Bucket([]byte(bucketPrefix + bucket))
	if b == nil {
		return nil, nil, nil, common.InvalidBucket(bucket)
	}
	// Retrieve mapping from keys to index for bucket
	m := tx.Bucket([]byte(bucket))
	if m == nil {
		return nil, nil, nil, fmt.Errorf("key-index mapping for bucket %s does not exist", bucket)
	}
	// Retrieve index
	idx, err := db.im.Get(bucket, db.kv)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("index for bucket %s does not exist, trace: %s", bucket, err.Error())
	}
	return b, m, idx, nil
}
