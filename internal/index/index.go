package index

type index interface {
	init() error
	close() error
	load() error
	save() error
	add() error
	delete() error
	search() error
}
