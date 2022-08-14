package priority_queue

import "sync"

type blockPriorityQueue[T any] struct {
	queue PriorityQueue[T]
	cond  *sync.Cond
}

func NewBlockPriorityQueue[T any]() PriorityQueue[T] {
	return &blockPriorityQueue[T]{NewBinaryHeap[T](), sync.NewCond(&sync.Mutex{})}
}

func (pq *blockPriorityQueue[T]) Size() int {
	return pq.queue.Size()
}

func (pq *blockPriorityQueue[T]) Push(x T, p int) {
	pq.cond.L.Lock()
	pq.queue.Push(x, p)
	pq.cond.Signal()
	pq.cond.L.Unlock()
}

func (pq *blockPriorityQueue[T]) Top() T {
	pq.cond.L.Lock()
	for pq.queue.Size() == 0 {
		pq.cond.Wait()
	}
	pq.cond.Signal()
	defer pq.cond.L.Unlock()
	return pq.queue.Top()
}

func (pq *blockPriorityQueue[T]) Pop() T {
	pq.cond.L.Lock()
	for pq.queue.Size() == 0 {
		pq.cond.Wait()
	}
	defer pq.cond.L.Unlock()
	return pq.queue.Pop()
}