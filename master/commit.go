package master

import (
	"heisenberg/model"
	"sync/atomic"
)

type committer struct {
	*Master
	pending   bool
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
	if c.committed.Load() != int32(c.manager.Size()) {
		return
	}
}

func (c *committer) start() {
	if c.pending || !c.wal.Full() {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.wal.Rotate()
	c.manager.Send(model.CommitRequest{})

	c.pending = true
}
