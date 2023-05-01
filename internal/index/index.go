package index

type index interface {
	Add() error
	Delete() error
	Search() error
}
