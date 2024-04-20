package store

import (
	"heisenberg/model"
	"heisenberg/state"
	"heisenberg/worker"
	"log/slog"
)

type Store struct {
	manager   worker.Manager
	compacter compacter
	disk      state.Disk
	heap      heap
	logger    *slog.Logger
}

func New() Store {
	return Store{}
}

func (s *Store) Receive(msg any) {
	switch msg := msg.(type) {
	case model.PutRequest:
		s.put(msg.Key, msg.Value)
	}
}

func (s *Store) put(key, value []byte) error {
	hp, err := s.heap.Write(value)
	if err != nil {
		return err
	}
	res := model.StorePutResponse{Key: key, Page: hp.Page, Offset: hp.Offset}
	s.manager.WhereJobEq(int(model.JobMaster)).Send(res)
	return nil
}

func process() {

}
