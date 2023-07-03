package doc

import (
	"heisenberg/utils"
	"regexp"
	"strings"
)

type strIndex interface {
	insert(id uint64, value string, old string) error
	remove(id uint64, value string) error
	search(value string) []uint64
}

type memStrIndex struct {
	index map[string][]uint64
}

func newMemStrIndex() *memStrIndex {
	return &memStrIndex{index: make(map[string][]uint64, 0)}
}

func (idx *memStrIndex) insert(id uint64, value string, old string) error {
	tokens := tokenise(value)
	oldTokens := tokenise(old)
	toAdd := utils.Except(tokens, oldTokens)
	toDelete := utils.Except(oldTokens, tokens)
	for _, token := range toDelete {
		idx.index[token] = utils.RemoveWhere(idx.index[token], id)
	}
	for _, token := range toAdd {
		idx.index[token] = append(idx.index[token], id)
	}
	return nil
}

func (idx *memStrIndex) remove(id uint64, value string) error {
	tokens := tokenise(value)
	for _, token := range tokens {
		idx.index[token] = utils.RemoveWhere(idx.index[token], id)
	}
	return nil
}

func (idx *memStrIndex) search(query string) []uint64 {
	queryTokens := tokenise(query)
	candidates := make([][]uint64, 0, 0)
	for _, token := range queryTokens {
		candidates = append(candidates, idx.index[token])
	}
	if len(candidates) == 0 {
		return nil
	}
	results := candidates[0]
	for _, candidate := range candidates[1:] {
		results = utils.Intersect(results, candidate)
	}
	return results
}

func tokenise(value string) []string {
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s]+")
	value = reg.ReplaceAllString(value, "")
	value = strings.ToLower(value)
	return strings.Fields(value)
}
