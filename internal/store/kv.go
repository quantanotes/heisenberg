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

func NewKv() (*kv, error) {
	return nil, nil
}

func LoadKv() (*kv, error) {
	return nil, nil
}

func CloseKv() {
	return
}

func (k *kv) Put(key []byte, value []byte, collection []byte) error {
	tx := func(tx *bbolt.Tx) error {
		b := tx.Bucket(collection)
		if b == nil {
			return internal.InvalidCollectionError(collection)
		}
		b.Put(key, value)
		return nil
	}
	return k.db.Update(tx)
}

func BeginTx() {

}
