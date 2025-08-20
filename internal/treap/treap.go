package treap

import (
	"errors"
	"fmt"
	"strings"
)
import "golang.org/x/exp/constraints"

var (
	ErrNilNode       = errors.New("treap: nil node")
	ErrNoParent      = errors.New("treap: node has no parent")
	ErrNotLeftChild  = errors.New("treap: node is not left child")
	ErrNotRightChild = errors.New("treap: node is not right child")
	ErrNotFound      = errors.New("treap: node not found")
)

type Treap[T constraints.Ordered] struct {
	root *Node[T]
}

type Node[T constraints.Ordered] struct {
	key      T
	priority float64
	left     *Node[T]
	right    *Node[T]
	parent   *Node[T]
}

func NewTreap[T constraints.Ordered]() *Treap[T] {
	return &Treap[T]{}
}

func NewNode[T constraints.Ordered](key T, priority float64) *Node[T] {
	return &Node[T]{key: key, priority: priority}
}

func (n *Node[T]) setLeft(node *Node[T]) {
	n.left = node
	if node != nil {
		node.parent = n
	}
}

func (n *Node[T]) setRight(node *Node[T]) {
	n.right = node
	if node != nil {
		node.parent = n
	}
}

func (t *Treap[T]) isRoot(x *Node[T]) bool { return t.root == x }

func (t *Treap[T]) rightRotate(x *Node[T]) error {
	if x == nil {
		return ErrNilNode
	}
	if t.isRoot(x) {
		return ErrNoParent
	}

	y := x.parent
	if y == nil {
		return ErrNoParent
	}
	if y.left != x {
		return ErrNotLeftChild
	}

	p := y.parent
	if t.isRoot(y) {
		t.root = x
		x.parent = nil
	} else if p.left == y {
		p.setLeft(x)
	} else {
		p.setRight(x)
	}

	y.setLeft(x.right)
	x.setRight(y)

	return nil
}

func (t *Treap[T]) leftRotate(x *Node[T]) error {
	if x == nil {
		return ErrNilNode
	}
	if t.isRoot(x) {
		return ErrNoParent
	}

	y := x.parent
	if y == nil {
		return ErrNoParent
	}
	if y.right != x {
		return ErrNotRightChild
	}

	p := y.parent
	if p != nil {
		if p.left == y {
			p.setLeft(x)
		} else {
			p.setRight(x)
		}
	} else {
		t.root = x
		x.parent = nil
	}

	y.setRight(x.left)
	x.setLeft(y)
	return nil
}

func (t *Treap[T]) Insert(key T, priority float64) error {
	newNode := NewNode(key, priority)

	if t.root == nil {
		t.root = newNode
		return nil
	}

	var parent *Node[T]
	node := t.root
	for node != nil {
		parent = node
		if key < node.key {
			node = node.left
		} else {
			node = node.right
		}
	}

	if key < parent.key {
		parent.setLeft(newNode)
	} else {
		parent.setRight(newNode)
	}

	for newNode.parent != nil && newNode.priority < newNode.parent.priority {
		if newNode == newNode.parent.left {
			if err := t.rightRotate(newNode); err != nil {
				return err
			}
		} else {
			if err := t.leftRotate(newNode); err != nil {
				return err
			}
		}
	}

	if newNode.parent == nil {
		t.root = newNode
	}
	return nil
}

func (t *Treap[T]) Remove(key T) bool {
	node := t.root
	for node != nil && node.key != key {
		if key < node.key {
			node = node.left
		} else {
			node = node.right
		}
	}
	if node == nil {
		return false
	}

	for !node.isLeaf() {
		if node.left != nil && (node.right == nil || node.left.priority <= node.right.priority) {
			if err := t.rightRotate(node.left); err != nil {
				return false
			}
		} else {
			if err := t.leftRotate(node.right); err != nil {
				return false
			}
		}
	}

	if node.parent == nil {
		t.root = nil
		return true
	}
	if node.parent.left == node {
		node.parent.left = nil
	} else {
		node.parent.right = nil
	}
	node.parent = nil

	return true
}

func (t *Treap[T]) Top() (T, error) {
	var zero T
	if t.root == nil {
		return zero, ErrNilNode
	}
	key := t.root.key
	if !t.Remove(key) {
		return zero, errors.New("treap: remove failed")
	}
	return key, nil
}

func (t *Treap[T]) Peek() (T, error) {
	var zero T
	if t.root == nil {
		return zero, ErrNilNode
	}
	return t.root.key, nil
}

func (t *Treap[T]) Update(key T, priority float64) error {
	if t.root == nil {
		return ErrNilNode
	}
	node := t.root.Search(key)
	if node == nil {
		return ErrNotFound
	}

	if node.priority == priority {
		return nil
	}

	old := node.priority
	node.priority = priority

	if priority < old {
		for node.parent != nil && node.priority < node.parent.priority {
			if node == node.parent.left {
				if err := t.rightRotate(node); err != nil {
					return err
				}
			} else {
				if err := t.leftRotate(node); err != nil {
					return err
				}
			}
		}
	} else {
		for {
			var leftOK = node.left != nil && node.left.priority < node.priority
			var rightOK = node.right != nil && node.right.priority < node.priority

			if leftOK && (!rightOK || node.left.priority <= node.right.priority) {
				if err := t.rightRotate(node.left); err != nil {
					return err
				}
			} else if rightOK {
				if err := t.leftRotate(node.right); err != nil {
					return err
				}
			} else {
				break
			}
		}
	}

	return nil
}

func (t *Treap[T]) Min() (T, error) {
	var zero T
	if t.root == nil {
		return zero, ErrNilNode
	}

	node := t.root
	for node.left != nil {
		node = node.left
	}

	return node.key, nil
}

func (t *Treap[T]) Max() (T, error) {
	var zero T
	if t.root == nil {
		return zero, ErrNilNode
	}

	node := t.root
	for node.right != nil {
		node = node.right
	}

	return node.key, nil
}

func (n *Node[T]) Search(targetKey T) *Node[T] {
	if n == nil {
		return nil
	}

	if n.key == targetKey {
		return n
	}

	if targetKey < n.key {
		return n.left.Search(targetKey)
	} else {
		return n.right.Search(targetKey)
	}
}

func (n *Node[T]) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func (t *Treap[T]) String() string {
	if t.root == nil {
		return "<empty treap>"
	}
	var b strings.Builder
	t.root.dump(&b, "", false)
	return b.String()
}

func (n *Node[T]) dump(b *strings.Builder, prefix string, isLeft bool) {
	if n == nil {
		return
	}
	if n.right != nil {
		n.right.dump(b, prefix+"    ", false)
	}

	if prefix == "" {
		fmt.Fprintf(b, "%v[p=%.3f]\n", n.key, n.priority) // корень
	} else if isLeft {
		fmt.Fprintf(b, "%s└── %v[p=%.3f]\n", prefix, n.key, n.priority)
	} else {
		fmt.Fprintf(b, "%s┌── %v[p=%.3f]\n", prefix, n.key, n.priority)
	}

	if n.left != nil {
		n.left.dump(b, prefix+"    ", true)
	}
}
