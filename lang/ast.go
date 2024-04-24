package lang

type Ast any

type From struct {
	Table string
}

type Select struct {
	Columns []string
}

type CreateTable struct {
	Name    string
	Columns []ColumnDefinition
}

type ColumnDefinition struct {
	Modifiers []ColumnModifier
	Type      Types
	Name      string
}

type ColumnModifier struct {
}

type Delete struct {
}

type FilterOp int

const (
	FilterOpEq FilterOp = iota
)

type Filter struct {
	Op      FilterOp
	Column  string
	Operand Ast
}
