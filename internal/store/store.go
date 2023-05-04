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

// func newStore(path string) (*store, error) {
// 	db, err := bbolt.Open(path, 0666, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &store{path, db, nil}, nil
// }

func loadStore(path string) (*store, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &store{path, db, nil}, nil
}

func (s *store) close() {
	s.db.Close()
}

func (s *store) createCollection(collection []byte) error {
	tx := func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket(collection)
		return err
	}
	return s.db.Update(tx)
}

func (s *store) deleteCollection(collection []byte) error {
	tx := func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket(collection)
		return err
	}
	return s.db.Update(tx)
}

func (s *store) get(key []byte, collection []byte) ([]byte, error) {
	var val []byte
	tx := func(tx *bbolt.Tx) error {
		b := tx.Bucket(collection)
		if b == nil {
			return internal.InvalidCollectionError(collection)
		}
		val = b.Get(key)
		if val == nil {
			return internal.InvalidKeyError(key, collection)
		}
		return nil
	}
	err := s.db.View(tx)
	if err != nil {
		return nil, err
	}
	return val, nil
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
