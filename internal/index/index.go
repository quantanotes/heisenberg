package index

import (
	"heisenberg/internal"
)

type index interface {
	init(space internal.SpaceType, dim int, max int) error
	close() error
	load() error
	save() error
	add() error
	delete() error
	search() error
}
