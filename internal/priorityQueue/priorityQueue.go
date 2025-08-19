package priorityQueue

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrElementNotFound = errors.New("element not found")
	ErrQueueIsEmpty    = errors.New("queue is empty")
)

func NewPriorityQueue[T comparable](d int, capacity int) *PriorityQueue[T] {
	if d < 2 {
		d = 2
	}
	if capacity < 0 {
		capacity = 0
	}
	return &PriorityQueue[T]{
		pairs:    make([]Pair[T], 0, capacity),
		sizeD:    d,
		indexMap: make(map[T]int, capacity),
	}
}

type PriorityQueue[T comparable] struct {
	pairs    []Pair[T]
	sizeD    int
	indexMap map[T]int
}

type Pair[T comparable] struct {
	priority float32
	value    T
}

func (q *PriorityQueue[T]) Top() (Pair[T], error) {
	if q.isEmpty() {
		return Pair[T]{}, ErrQueueIsEmpty
	}

	p := q.removeLast()

	if q.isEmpty() {
		return p, nil
	}

	element := q.pairs[0]
	q.pairs[0] = p
	q.indexMap[p.value] = 0
	q.pushDown()
	delete(q.indexMap, element.value)
	return element, nil
}

func (q *PriorityQueue[T]) Peek() (Pair[T], error) {
	if q.isEmpty() {
		return Pair[T]{}, ErrQueueIsEmpty
	}

	return q.pairs[0], nil
}

func (q *PriorityQueue[T]) Insert(element T, priority float32) {
	newPair := Pair[T]{value: element, priority: priority}
	q.pairs = append(q.pairs, newPair)
	q.indexMap[element] = len(q.pairs) - 1
	q.bubbleUp()
}

func (q *PriorityQueue[T]) Remove(element T) error {
	index, ok := q.indexMap[element]
	if !ok {
		return ErrElementNotFound
	}

	lastIndex := len(q.pairs) - 1
	if index == lastIndex {
		q.pairs = q.pairs[:lastIndex]
		delete(q.indexMap, element)
		return nil
	}

	removedElement := q.pairs[index].value
	removedPriority := q.pairs[index].priority

	q.pairs[index] = q.pairs[lastIndex]
	q.indexMap[q.pairs[index].value] = index

	q.pairs = q.pairs[:lastIndex]
	delete(q.indexMap, removedElement)

	if q.pairs[index].priority < removedPriority {
		q.bubbleUpIndex(index)
	} else if q.pairs[index].priority > removedPriority {
		q.pushDownIndex(index)
	}

	return nil
}

func (q *PriorityQueue[T]) Update(element T, newPriority float32) error {
	index, ok := q.indexMap[element]
	if !ok {
		return ErrElementNotFound
	}

	oldPriority := q.pairs[index].priority
	q.pairs[index].priority = newPriority

	if newPriority < oldPriority {
		q.bubbleUpIndex(index)
	} else if newPriority > oldPriority {
		q.pushDownIndex(index)
	}

	return nil
}

func (q *PriorityQueue[T]) heapify() {
	q.indexMap = make(map[T]int, len(q.pairs))
	for i, pair := range q.pairs {
		q.indexMap[pair.value] = i
	}

	for index := (len(q.pairs) - 1) / q.sizeD; index >= 0; index-- {
		q.pushDownIndex(index)
	}
}

func (q *PriorityQueue[T]) bubbleUp() {
	q.bubbleUpIndex(len(q.pairs) - 1)
}

func (q *PriorityQueue[T]) bubbleUpIndex(index int) {
	current := q.pairs[index]
	for index > 0 {
		parentIndex := q.getParentIndex(index)
		if q.pairs[parentIndex].priority < current.priority {
			q.pairs[index] = q.pairs[parentIndex]
			q.indexMap[q.pairs[index].value] = index
			index = parentIndex
		} else {
			break
		}
	}
	q.pairs[index] = current
	q.indexMap[current.value] = index
}

func (q *PriorityQueue[T]) pushDown() {
	q.pushDownIndex(0)
}

func (q *PriorityQueue[T]) pushDownIndex(currentIndex int) {
	for currentIndex < q.firstLeafIndex() {
		_, childIndex := q.highestPriorityChild(currentIndex)
		if childIndex == -1 {
			break
		}
		if q.pairs[childIndex].priority > q.pairs[currentIndex].priority {
			q.pairs[currentIndex], q.pairs[childIndex] = q.pairs[childIndex], q.pairs[currentIndex]
			q.indexMap[q.pairs[currentIndex].value] = currentIndex
			q.indexMap[q.pairs[childIndex].value] = childIndex
			currentIndex = childIndex
		} else {
			break
		}
	}
}

func (q *PriorityQueue[T]) getParentIndex(parentIndex int) int {
	return (parentIndex - 1) / q.sizeD
}

func (q *PriorityQueue[T]) firstLeafIndex() int {
	return (len(q.pairs)-2)/q.sizeD + 1
}

func (q *PriorityQueue[T]) isEmpty() bool {
	return len(q.pairs) == 0
}

func (q *PriorityQueue[T]) removeLast() Pair[T] {
	element := q.pairs[len(q.pairs)-1]
	delete(q.indexMap, element.value)
	q.pairs = q.pairs[:len(q.pairs)-1]
	return element
}

func (q *PriorityQueue[T]) highestPriorityChild(currentIndex int) (best Pair[T], bestIdx int) {
	start := currentIndex*q.sizeD + 1
	if start >= len(q.pairs) {
		return Pair[T]{}, -1
	}
	end := start + q.sizeD - 1
	if end >= len(q.pairs) {
		end = len(q.pairs) - 1
	}

	bestIdx = start
	best = q.pairs[start]
	for i := start + 1; i <= end; i++ {
		if q.pairs[i].priority > best.priority {
			best = q.pairs[i]
			bestIdx = i
		}
	}
	return best, bestIdx
}

func (q *PriorityQueue[T]) AsciiTree() string {
	var b strings.Builder
	n := len(q.pairs)
	if n == 0 {
		return "(empty)\n"
	}
	b.WriteString(q.nodeLabel(0) + "\n")
	start := 0*q.sizeD + 1
	end := start + q.sizeD
	if end > n {
		end = n
	}
	for i := start; i < end; i++ {
		last := i == end-1
		q.asciiTree(i, "", last, &b)
	}
	return b.String()
}

func (q *PriorityQueue[T]) asciiTree(i int, prefix string, isLast bool, b *strings.Builder) {
	connector := "├── "
	childPrefix := prefix + "│   "
	if isLast {
		connector = "└── "
		childPrefix = prefix + "    "
	}
	b.WriteString(prefix + connector + q.nodeLabel(i) + "\n")

	n := len(q.pairs)
	start := i*q.sizeD + 1
	if start >= n {
		return
	}
	end := start + q.sizeD
	if end > n {
		end = n
	}
	for j := start; j < end; j++ {
		last := j == end-1
		q.asciiTree(j, childPrefix, last, b)
	}
}

func (q *PriorityQueue[T]) nodeLabel(i int) string {
	p := q.pairs[i]
	return fmt.Sprintf("[%.1f] %v", p.priority, p.value)
}
