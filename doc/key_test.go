package doc

import (
	"fmt"
	"testing"
)

func TestKey(t *testing.T) {
	k := key{"hello", "world"}
	compressed := k.Compress()
	fmt.Println(compressed, len(compressed))
	uncompressed := UncompressKey(compressed)
	for _, s := range uncompressed {
		fmt.Println(s, len(s))
	}
}
