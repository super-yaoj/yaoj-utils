package pq

//-------Blocked priority queue

type queueEntry[T any] struct {
	val      T
	priority int
}

type PriorityQueue[T any] interface {
	Size() int
	Push(T, int)
	Top() T
	Pop() T
}

type binaryHeap[T any] []queueEntry[T]

func NewBinaryHeap[T any]() PriorityQueue[T] {
	return &binaryHeap[T]{}
}

func (bh *binaryHeap[T]) Size() int {
	return len(*bh)
}

func (bh *binaryHeap[T]) Push(x T, p int) {
	*bh = append(*bh, queueEntry[T]{})
	n := len(*bh) - 1
	for n > 0 {
		f := (n - 1) / 2
		if (*bh)[f].priority < p {
			(*bh)[n] = (*bh)[f]
		} else {
			break
		}
		n = f
	}
	(*bh)[n] = queueEntry[T]{x, p}
}

func (bh *binaryHeap[T]) Pop() T {
	top := (*bh)[0]
	n := len(*bh) - 1
	(*bh)[0], (*bh)[n] = (*bh)[n], (*bh)[0]
	(*bh) = (*bh)[:n]
	if n > 0 {
		temp := (*bh)[0]
		i := 0
		for {
			l := i*2 + 1
			r := l + 1
			if l >= n {
				break
			}
			if r >= n {
				if (*bh)[l].priority > temp.priority {
					(*bh)[i] = (*bh)[l]
					i = l
				} else {
					break
				}
			} else if (*bh)[l].priority < (*bh)[r].priority {
				if (*bh)[r].priority > temp.priority {
					(*bh)[i] = (*bh)[r]
					i = r
				} else {
					break
				}
			} else {
				if (*bh)[l].priority > temp.priority {
					(*bh)[i] = (*bh)[l]
					i = l
				} else {
					break
				}
			}
		}
		(*bh)[i] = temp
	}
	return top.val
}

func (bh *binaryHeap[T]) Top() T {
	return (*bh)[0].val
}