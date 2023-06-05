package core

import (
	"fmt"

	"github.com/quantanotes/heisenberg/common"
	"github.com/quantanotes/heisenberg/core/hnsw"
)

type IndexerType int

const (
	HNSWIndexerType  IndexerType = 1
	AnnoyIndexerType IndexerType = 2
)

type Index interface {
	Close()
	Save(string) error
	Insert(uint, []float32) error
	Delete(uint) error
	Search([]float32, uint) ([]uint, error)
	Next() uint
	GetConfig() common.IndexConfig
}

func NewIndex(indexer IndexerType, name string, dim uint, space common.SpaceType) (Index, error) {
	config := common.IndexConfig{
		Name:     name,
		Indexer:  int(indexer),
		FreeList: make([]uint, 0),
		Dim:      dim,
		Space:    int(space),
		Count:    0,
	}

	switch indexer {
	case HNSWIndexerType:
		return hnsw.New(space, config, hnsw.HNSWOptions{M: 50, Ef: 100, Max: 1000000}), nil
	case AnnoyIndexerType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown indexer type")
	}
}

func LoadIndex(path string, config common.IndexConfig) (Index, error) {
	switch IndexerType(config.Indexer) {
	case HNSWIndexerType:
		return hnsw.Load(path, config)
	case AnnoyIndexerType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown indexer type")
	}
}
