package store

import (
	"heisenberg/ioutil"
	"heisenberg/state"
)

type heap struct {
	pager       state.Pager
	headSize    int
	maxPageSize int
}

type heapPointer struct {
	Page   uint32
	Offset uint32
}

func (h *heap) Write(value []byte) (heapPointer, error) {
	h.checkRotateHead()
	page := uint32(h.pager.HeadIndex())
	offset := uint32(h.headSize)
	head := h.pager.Head()
	if err := ioutil.AppendBytes(head, head, value); err != nil {
		return heapPointer{}, err
	}
	return heapPointer{page, offset}, nil
}

func (h *heap) checkRotateHead() error {
	if h.headSize > h.maxPageSize {
		if _, err := h.pager.New(); err != nil {
			return err
		}
		h.headSize = 0
	}
	return nil
}
