package master

import (
	"heisenberg/model"
	"sync/atomic"
)

// Handles rotating WAL/Memtable during commit phase.
type committer struct {
	*Master
	running   bool
	committed atomic.Int32
}

func (c *committer) Process() {
	c.start()
	c.finish()
}

func (c *committer) Commit() {
	c.committed.Add(1)
}

func (c *committer) finish() {
	// If all storage nodes have responded, we can finish the rotation
	if c.committed.Load() != int32(c.manager.Size()) {
		c.committed.Store(0)
		c.running = false
		return
	}
}

func (c *committer) start() {
	// Begin the commit phase once the WAL is full
	if c.running || !c.wal.Full() {
		return
	}

	c.running = true
	c.wal.Rotate()
	c.memtable.Rotate()
	c.manager.Send(model.CommitRequest{})
}
