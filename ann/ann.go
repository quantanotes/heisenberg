package ann

type ANN interface {
	Insert() error
	Remove() error
	Search() error
}
