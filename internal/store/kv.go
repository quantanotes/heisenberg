package store

import (
	"heisenberg/internal"

	"go.etcd.io/bbolt"
)

// Base key-value storage
type kv struct {
	path   string               // Storage location
	db     *bbolt.DB            // BoltDB storage backend
	txPool map[string]*bbolt.Tx // Pending transactions
}

func newKv() (*kv, error) {
	return nil, nil
}

func loadKv() (*kv, error) {
	return nil, nil
}

func (kv *kv) closeKv() {
	kv.db.Close()
}

func (kv *kv) get(key []byte, collection []byte) ([]byte, error) {
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
	return value, kv.db.View(tx)
}

func (kv *kv) put(key []byte, value []byte, collection []byte) error {
	tx := func(tx *bbolt.Tx) error {
		b := tx.Bucket(collection)
		if b == nil {
			return internal.InvalidCollectionError(collection)
		}
		b.Put(key, value)
		return nil
	}
	return kv.db.Update(tx)
}

func (kv *kv) delete(key []byte, collection []byte) error {
	tx := func(tx *bbolt.Tx) error {
		b := tx.Bucket(collection)
		if b == nil {
			return internal.InvalidCollectionError(collection)
		}
		return b.Delete(key)
	}
	return kv.db.Update(tx)
}

func beginTx() {

}
