package priorityQueueByLinkedList

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

var (
	ErrEmpty    = errors.New("priorityQueue: is empty")
	ErrNotFound = errors.New("priorityQueue: node not found")
)

type PriorityQueue[T constraints.Ordered] struct {
	root  *Node[T]
	sizeD int
	size  int
}

type Node[T constraints.Ordered] struct {
	childes []*Node[T]
	parent  *Node[T]
	pair    Pair[T]
}

type Pair[T constraints.Ordered] struct {
	priority float64
	value    T
}

func NewPriorityQueue[T constraints.Ordered](d int) *PriorityQueue[T] {
	if d < 2 {
		d = 2
	}
	return &PriorityQueue[T]{sizeD: d}
}

func (p *PriorityQueue[T]) Top() (Pair[T], error) {
	if p.root == nil {
		var zero Pair[T]
		return zero, ErrEmpty
	}
	maxPair := p.root.pair

	if p.size == 1 {
		p.root = nil
		p.size = 0
		return maxPair, nil
	}

	last := p.nodeAt(p.size)
	p.root.pair = last.pair
	p.detachLast()
	p.pushDownIndex(1)

	return maxPair, nil
}

func (p *PriorityQueue[T]) Peek() (Pair[T], error) {
	if p.root == nil {
		var zero Pair[T]
		return zero, ErrEmpty
	}
	return p.root.pair, nil
}

func (p *PriorityQueue[T]) Insert(element T, priority float64) {
	newPair := Pair[T]{priority: priority, value: element}

	p.size++
	if p.root == nil {
		p.root = &Node[T]{childes: make([]*Node[T], p.sizeD), pair: newPair}
		return
	}

	parent, childIdx := p.insertionParentAndIdx(p.size)
	node := &Node[T]{childes: make([]*Node[T], p.sizeD), parent: parent, pair: newPair}
	parent.childes[childIdx] = node
	p.bubbleUpNode(node)
}

func (p *PriorityQueue[T]) Remove(element T) error {
	if p.size == 0 {
		return ErrEmpty
	}

	var target *Node[T]
	var targetIdx int
	for i := 1; i <= p.size; i++ {
		n := p.nodeAt(i)
		if n != nil && n.pair.value == element {
			target = n
			targetIdx = i
			break
		}
	}
	if target == nil {
		return ErrNotFound
	}

	if p.size == 1 {
		p.root = nil
		p.size = 0
		return nil
	}

	if targetIdx == p.size {
		p.detachLast()
		return nil
	}

	last := p.nodeAt(p.size)
	oldPriority := target.pair.priority
	target.pair = last.pair
	p.detachLast()

	if target.pair.priority > oldPriority {
		p.bubbleUpNode(target)
	} else {
		p.pushDownNode(target)
	}
	return nil
}

func (p *PriorityQueue[T]) Update(element T, newPriority float64) error {
	if p.size == 0 {
		return ErrEmpty
	}
	var node *Node[T]
	for i := 1; i <= p.size; i++ {
		n := p.nodeAt(i)
		if n != nil && n.pair.value == element {
			node = n
			break
		}
	}
	if node == nil {
		return ErrNotFound
	}
	old := node.pair.priority
	node.pair.priority = newPriority
	if newPriority > old {
		p.bubbleUpNode(node)
	} else if newPriority < old {
		p.pushDownNode(node)
	}
	return nil
}

func (p *PriorityQueue[T]) heapify() {
	if p.size <= 1 {
		return
	}
	start := ((p.size - 2) / p.sizeD) + 1
	for i := start; i >= 1; i-- {
		p.pushDownIndex(i)
	}
}

func (p *PriorityQueue[T]) bubbleUp() {
	if p.size == 0 {
		return
	}
	p.bubbleUpIndex(p.size)
}

func (p *PriorityQueue[T]) bubbleUpIndex(index int) {
	n := p.nodeAt(index)
	if n == nil {
		return
	}
	p.bubbleUpNode(n)
}

func (p *PriorityQueue[T]) bubbleUpNode(n *Node[T]) {
	for n.parent != nil && n.pair.priority > n.parent.pair.priority {
		n.pair, n.parent.pair = n.parent.pair, n.pair
		n = n.parent
	}
}

func (q *PriorityQueue[T]) pushDown() {
	q.pushDownIndex(1)
}

func (q *PriorityQueue[T]) pushDownIndex(currentIndex int) {
	n := q.nodeAt(currentIndex)
	if n == nil {
		return
	}
	q.pushDownNode(n)
}

func (q *PriorityQueue[T]) pushDownNode(n *Node[T]) {
	for {
		var best *Node[T]
		bestIdx := -1
		for i, ch := range n.childes {
			if ch == nil {
				continue
			}
			if best == nil || ch.pair.priority > best.pair.priority {
				best = ch
				bestIdx = i
			}
		}
		if best == nil || best.pair.priority <= n.pair.priority {
			return
		}
		_ = bestIdx
		n.pair, best.pair = best.pair, n.pair
		n = best
	}
}

func (p *PriorityQueue[T]) pathTo(index int) []int {
	if index <= 1 {
		return nil
	}
	m := index - 1
	path := make([]int, 0, 8)
	for m > 0 {
		child := (m - 1) % p.sizeD
		path = append(path, child)
		m = (m - 1) / p.sizeD
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func (p *PriorityQueue[T]) nodeAt(index int) *Node[T] {
	if p.root == nil || index < 1 || index > p.size {
		return nil
	}
	n := p.root
	for _, ci := range p.pathTo(index) {
		if n == nil {
			return nil
		}
		n = n.childes[ci]
	}
	return n
}

func (p *PriorityQueue[T]) insertionParentAndIdx(index int) (*Node[T], int) {
	path := p.pathTo(index)
	childIdx := path[len(path)-1]
	parentPath := path[:len(path)-1]
	parent := p.root
	for _, ci := range parentPath {
		parent = parent.childes[ci]
	}
	return parent, childIdx
}

func (p *PriorityQueue[T]) detachLast() {
	if p.size == 0 {
		return
	}
	if p.size == 1 {
		p.root = nil
		p.size = 0
		return
	}
	parent, childIdx := p.insertionParentAndIdx(p.size)
	parent.childes[childIdx] = nil
	p.size--
}

func (p *PriorityQueue[T]) AsciiTree() string {
	var b strings.Builder
	n := p.size
	if n == 0 {
		return "(empty)\n"
	}
	b.WriteString(p.nodeLabel(0) + "\n")

	start := 0*p.sizeD + 1
	end := start + p.sizeD
	if end > n {
		end = n
	}
	for i := start; i < end; i++ {
		last := i == end-1
		p.asciiTree(i, "", last, &b)
	}
	return b.String()
}

func (p *PriorityQueue[T]) asciiTree(i int, prefix string, isLast bool, b *strings.Builder) {
	connector := "├── "
	childPrefix := prefix + "│   "
	if isLast {
		connector = "└── "
		childPrefix = prefix + "    "
	}
	b.WriteString(prefix + connector + p.nodeLabel(i) + "\n")

	n := p.size
	start := i*p.sizeD + 1
	if start >= n {
		return
	}
	end := start + p.sizeD
	if end > n {
		end = n
	}
	for j := start; j < end; j++ {
		last := j == end-1
		p.asciiTree(j, childPrefix, last, b)
	}
}

func (p *PriorityQueue[T]) nodeLabel(i int) string {
	n := p.nodeAt(i + 1)
	if n == nil {
		return fmt.Sprintf("(nil #%d)", i)
	}
	return fmt.Sprintf("[%.1f] %v", n.pair.priority, n.pair.value)
}
