package core

import (
	"fmt"
)

type IndexerType int

const (
	HNSWIndexerType  IndexerType = 1
	AnnoyIndexerType             = 2
)

type IndexConfig struct {
	Name     string
	FreeList []uint
	Dim      uint
	Space    uint
}

type Index interface {
	Close()
	Save(string) error
	Insert(uint, []float32) error
	Delete(uint) error
	Search([]float32, uint) ([]int, error)
	Next() uint
	GetConfig() IndexConfig
}


func NewIndex(config IndexConfig, indexer IndexerType) (*Index, error) {
	switch indexer {
	case HNSWIndexerType:
		return nil, nil
	case AnnoyIndexerType:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown indexer type")
	}
}

func LoadIndex(path string, config IndexConfig) (*Index, error) {
	return nil, nil
}
