package master

import (
	"heisenberg/model"
	"heisenberg/query"
	"heisenberg/worker"
	"sync"
	"sync/atomic"
)

type Master struct {
	mu        sync.Mutex
	stop      chan struct{}
	manager   *worker.Manager
	wal       wal
	memtable  memtable
	committer committer
}

func New() *Master {
	m := &Master{}
	m.committer = committer{m, false, atomic.Int32{}}
	go m.process()
	return m
}

func (m *Master) Receive(msg any) {
	switch msg.(type) {
	case model.CommitRequest:
		m.manager.Send(model.CommitResponse{})
	case model.CommitResponse:
		m.committer.Commit()
	}
}

func (m *Master) Stop() {
	m.stop <- struct{}{}
	m.wal.Close()
}

func (m *Master) process() {
	for {
		select {
		case <-m.stop:
			return
		default:
			m.committer.Process()
		}
	}
}

func (m *Master) put(key, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.wal.Write(model.MethodPut, key, value); err != nil {
		return err
	}

	m.memtable.Put(key, value)

	return nil
}

func (m *Master) execute(plan query.Plan) {

}
