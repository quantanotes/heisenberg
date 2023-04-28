package store

import (
	"heisenberg/internal"

	"go.etcd.io/bbolt"
)

// Base key-value storage
type store struct {
	path   string               // Storage location
	db     *bbolt.DB            // BoltDB storage backend
	txPool map[string]*bbolt.Tx // Pending transactions
}

func newKv() (*store, error) {
	return nil, nil
}

func loadKv() (*store, error) {
	return nil, nil
}

func (s *store) closeKv() {
	s.db.Close()
}

func (s *store) get(key []byte, collection []byte) ([]byte, error) {
	var value []byte
	tx := func(tx *bbolt.Tx) error {
		b := tx.Bucket(collection)
		if b == nil {
			return internal.InvalidCollectionError(collection)
		}
		value = b.Get(key)
		if value == nil {
			return internal.InvalidKeyError(key, collection)
		}
		return nil
	}
	return value, s.db.View(tx)
}

func (s *store) put(key []byte, value []byte, collection []byte) error {
	tx := func(tx *bbolt.Tx) error {
		b := tx.Bucket(collection)
		if b == nil {
			return internal.InvalidCollectionError(collection)
		}
		b.Put(key, value)
		return nil
	}
	return s.db.Update(tx)
}

func (s *store) delete(key []byte, collection []byte) error {
	tx := func(tx *bbolt.Tx) error {
		b := tx.Bucket(collection)
		if b == nil {
			return internal.InvalidCollectionError(collection)
		}
		return b.Delete(key)
	}
	return s.db.Update(tx)
}

func beginTx() {

}
