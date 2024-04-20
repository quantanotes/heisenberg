package query

type nodeType uint32

const (
	NodeRoot nodeType = iota
	NodeScan
	NodeSeek
	NodeGet
	NodePut
	NodeDelete
	NodeFilter
	NodeApply
	NodeSort
)

type node struct {
	Type     nodeType
	Detail   Detail
	Children []uint32
}

func newNode(t nodeType, detail Detail) node {
	return node{
		Type:     t,
		Detail:   detail,
		Children: []uint32{},
	}
}
