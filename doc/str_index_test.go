package doc

import (
	"heisenberg/utils"
	"testing"
)

func TestStrIndex(t *testing.T) {
	index := newMemStrIndex()

	index.insert(1, "apple banana", "")
	index.insert(2, "banana cherry", "")
	index.insert(3, "apple cherry", "")
	index.insert(4, "apple banana cherry", "")

	results := index.search("apple")
	expected := []uint64{1, 3, 4}
	if !utils.Equals(results, expected) {
		t.Errorf("Search results for 'apple' are incorrect. Got %v, expected %v", results, expected)
	}

	index.remove(2, "banana cherry")

	results = index.search("banana")
	expected = []uint64{1, 4}
	if !utils.Equals(results, expected) {
		t.Errorf("Search results for 'banana' are incorrect. Got %v, expected %v", results, expected)
	}
}
