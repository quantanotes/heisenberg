package store

type Store interface {
	Get() error
	Put() error
}

type Tx interface {
	Commit() error
	Rollback() error
}
