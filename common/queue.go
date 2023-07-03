package common

type Queue[T any] struct {
	data []T
}

func NewQueue[T any]() Queue[T] {
	return Queue[T]{data: make([]T, 0)}
}

func (q *Queue[T]) Push(item T) {
	q.data = append(q.data, item)
}

func (q *Queue[T]) Pop() *T {
	if len(q.data) == 0 {
		return nil
	}
	item := q.data[0]
	q.data = q.data[1:len(q.data)]
	return &item
}

func (q *Queue[T]) Peek() *T {
	if len(q.data) == 0 {
		return nil
	}
	return &q.data[0]
}
