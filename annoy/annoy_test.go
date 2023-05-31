package annoy

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/quantanotes/heisenberg/annoy/annoyindex"
)

func TestAnnoyIndex(t *testing.T) {
	f := 40
	idx := annoyindex.NewAnnoyIndexAngular(f)
	for i := 0; i < 1000; i++ {
		item := make([]float32, 0, f)
		for x := 0; x < f; x++ {
			item = append(item, rand.Float32())
		}
		idx.AddItem(i, item)
	}
	idx.Build(10)
	idx.Save("test.ann")

	annoyindex.DeleteAnnoyIndexAngular(idx)

	idx = annoyindex.NewAnnoyIndexAngular(f)
	idx.Load("test.ann")

	var result []int
	idx.GetNnsByItem(0, 1000, -1, &result)
	fmt.Printf("%v\n", result)
}
