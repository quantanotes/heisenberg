package plan

type NodeKind int

const (
	NodePut = iota
)

type Node struct {
	Children []int
}

type Plan struct {
	Nodes []Node
}
