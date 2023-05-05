package index

import "sort"

type distqueue[T interface{}] struct {
	items []T
	size  int
	dist  func(T) float32
}

func newDistqueue[T interface{}](size int, dist func(T) float32) *distqueue[T] {
	return &distqueue[T]{
		items: make([]T, 0),
		size:  size,
		dist:  dist,
	}
}

func (q *distqueue[T]) less(i, j int) bool {
	return q.dist(q.items[i]) < q.dist(q.items[j])
}

func (q *distqueue[T]) push(item T) {
	q.items = append(q.items, item)
	sort.Slice(q.items, q.less)
	q.items = q.items[:q.size]
}

func (q *distqueue[T]) pushMany(items []T) {
	q.items = append(q.items, items...)
	sort.Slice(q.items, q.less)
	q.items = q.items[:q.size]
}
