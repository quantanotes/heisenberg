package graph

type EdgeList = []*Edge

type Edge struct {
	from string
	to   string
	meta map[string]any
}
