package master

import (
	"heisenberg/model"
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
	// Pending responses from asynchronous manager.Send()
	// This needs to be 0 before the committer can active
	pending atomic.Int32
}

func New() *Master {
	m := &Master{}
	m.committer = committer{m, false, atomic.Int32{}}
	go m.process()
	return m
}

func (m *Master) Receive(msg any) {
	switch msg := msg.(type) {
	case model.StorePutResponse:
		m.memtable.UpdateHeapPointer(msg.Key, msg.Page, msg.Offset)
		m.pending.Add(-1)
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
	// Signal to store nodes to append to the value log
	m.pending.Add(1)
	m.manager.
		WhereJobEq(int(model.JobStore)).
		Send(model.PutRequest{Key: key, Value: value})

	return nil
}

// func (m *Master) execute(plan query.Plan) {

// }
