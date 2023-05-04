package index

type spaceType = int

const (
	ip     spaceType = 1
	cosine           = 2
	l2               = 3
)

type index interface {
	init(space spaceType, dim int, max int) error
	close() error
	load() error
	save() error
	add() error
	delete() error
	search() error
}
