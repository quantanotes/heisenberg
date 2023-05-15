package core

import (
	"fmt"
	"heisenberg/utils"
	"os"
	"path/filepath"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"go.etcd.io/bbolt"
)

type indexMap = *orderedmap.OrderedMap[string, *Index]

// Manages the memory usage of multiple indices.
type IndexManager struct {
	indices indexMap
	max     uint // Maximum number of indices in memory
	path    string
}

func (im *IndexManager) New(conf indexConfig) {
	idx := NewIndex(conf)
	idx.Save(im.GetPath(conf.Name))
	im.push(idx)
}

func NewIndexManager(path string, max uint) *IndexManager {
	return &IndexManager{
		indices: orderedmap.New[string, *Index](),
		max:     max,
		path:    path,
	}
}

func (im *IndexManager) Close() {
	for pair := im.indices.Oldest(); pair != nil; pair = pair.Next() {
		pair.Value.Save(im.GetPath(pair.Value.config.Name))
		pair.Value.Close()
	}
}

func (im *IndexManager) Get(name string, db *bbolt.DB) (*Index, error) {
	idx, ok := im.indices.Get(name)
	if !ok {
		fmt.Printf("loading %s", name)
		return im.load(name, db)
	}
	im.indices.MoveToBack(name)
	return idx, nil
}

func (im *IndexManager) Delete(name string) error {
	return os.Remove(im.GetPath(name))
}

func (im *IndexManager) load(name string, db *bbolt.DB) (*Index, error) {
	// Retrieve configuration from key value store
	conf := &indexConfig{}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(name))
		data := b.Get([]byte(configKey))
		if data == nil {
			return utils.InvalidIndexConfig(name)
		}
		err := utils.FromBytes(data, conf)
		if err != nil {
			return utils.InvalidIndexConfig(name, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Load index from disk
	idx, err := LoadIndex(im.GetPath(name), *conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Push index to queue
	im.push(idx)
	return idx, nil
}

func (im *IndexManager) push(idx *Index) {
	im.indices.Store(idx.config.Name, idx)
	for im.indices.Len() > int(im.max) {
		pair := im.indices.Oldest()
		pair.Value.Save(im.GetPath(idx.config.Name))
		im.indices.Delete(pair.Key)
	}
}

func (im *IndexManager) GetPath(name string) string {
	return filepath.Join(im.path, "/"+name+".idx")
}
