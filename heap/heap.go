package heap

import (
	"heisenberg/ioutil"
	"heisenberg/state"
)

type Heap struct {
	pager       state.Pager
	headSize    int
	maxPageSize int
}

type Pointer struct {
	Page   uint32
	Offset uint32
}

func (h *Heap) write(value []byte) (Pointer, error) {
	if len(value)+h.headSize > h.maxPageSize {
		if _, err := h.pager.New(); err != nil {
			return Pointer{}, err
		}
	}

	page := uint32(h.pager.HeadIndex())
	offset := uint32(h.headSize)

	if err := ioutil.AppendBytes(h.pager.Head(), h.pager.Head(), value); err != nil {
		return Pointer{}, err
	}

	return Pointer{
		page,
		offset,
	}, nil
}
