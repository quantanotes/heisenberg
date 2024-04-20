package db

import (
	"heisenberg/query"
	"heisenberg/worker"
)

type DB struct {
	manager worker.Manager
}

func New() *DB {
	return &DB{}
}

func (db *DB) get() {

}

func (db *DB) put() {
}

func (db *DB) delete() {

}

func (db *DB) execute(plan query.Plan) {

}
