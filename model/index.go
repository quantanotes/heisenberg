package model

const (
	IndexBlockSize = 256
)

type Index struct {
	Header IndexHeader
	Blocks []IndexBlock
}

type IndexBlock struct {
	Entries []IndexEntry
}

type IndexEntry struct {
	Key    []byte
	Page   uint32 // heap pointer page
	Offset uint32 // heap pointer offset
}

type IndexHeader struct {
	First  []byte // first key in index
	Last   []byte // last key in index
	Blocks []IndexBlockHeader
}

type IndexBlockHeader struct {
	First  []byte // first key in block
	Offset uint32 // block offset in page
}
