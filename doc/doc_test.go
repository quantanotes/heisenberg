package doc

import (
	"fmt"
	"heisenberg/common"
	"testing"
)

func TestDocFlatten(t *testing.T) {
	data := common.Meta{
		"foo": "bar",
		"foobar": common.Meta{
			"baz": "quz",
			"foobaz": common.Meta{
				"quz": "quxbar",
			},
		},
	}
	flattened := flattenDoc(data)
	fmt.Printf("%v\n", flattened)
	unflattened := unflattenDoc(flattened)
	fmt.Printf("%v\n", unflattened)
}
