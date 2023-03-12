// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

type Comparison[T any] func(T, T) int

type builtin interface {
	~string | ~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8
}

func Compare[C builtin](x, y C) int {
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	default:
		return 0
	}
}

// TreeSet provides a sorted set implementation.
//
// The underlying data-structure is a standard Red-Black Tree.
// https://en.wikipedia.org/wiki/Red–black_tree
//
// The implementation prioritizes readability over maximal optimizations.
type TreeSet[S any, C Comparison[S]] struct {
	comparison C
	root       *node[S]
	size       int
}

func NewTreeSet[T any, C Comparison[T]](compare C) *TreeSet[T, C] {
	return &TreeSet[T, C]{
		comparison: compare,
		root:       nil,
		size:       0,
	}
}

// Insert item into s.
//
// Returns true if s was modified (item was not already in s), false otherwise.
func (s *TreeSet[T, C]) Insert(item T) bool {
	return s.insert(&node[T]{
		element: item,
		color:   red,
	})
}

// Min returns the smallest item in the set.
//
// Must not be called on an empty set.
func (s *TreeSet[T, C]) Min() T {
	if s.root == nil {
		panic("min: tree is empty")
	}
	n := s.min(s.root)
	return n.element
}

// Max returns the largest item in s.
//
// Must not be called on an empty set.
func (s *TreeSet[T, C]) Max() T {
	if s.root == nil {
		panic("max: tree is empty")
	}
	n := s.max(s.root)
	return n.element
}

// Size returns the number of elements in s.
func (s *TreeSet[T, C]) Size() int {
	return s.size
}

// Empty returns true if there are no elements in s.
func (s *TreeSet[T, C]) Empty() bool {
	return s.Size() == 0
}

// Red-Black Tree Invariants
//
// 1. each node is either red or black
// 2. the root node is always black
// 3. nil leaf nodes are always black
// 4. a red node must not have red children
// 5. all simple paths from a node to nil leaf contain the same number of
// black nodes

type color bool

const (
	red   color = false
	black color = true
)

type node[T any] struct {
	element T
	color   color
	parent  *node[T]
	left    *node[T]
	right   *node[T]
}

func (n *node[T]) less(c Comparison[T], o *node[T]) bool {
	return c(n.element, o.element) < 0
}

func (n *node[T]) greater(c Comparison[T], o *node[T]) bool {
	return c(n.element, o.element) > 0
}

func (n *node[T]) black() bool {
	return n.color == black
}

func (n *node[T]) red() bool {
	return n.color == red
}

// func (t *TreeSet[S, C]) locate(n *node[S], target S) *node[S] {
// 	if n == nil || n.element.Equal(target) {
// 		return n
// 	}

// 	if n.element.Less(target) {
// 		return t.locate(n.right, target)
// 	}

// 	return t.locate(n.left, target)
// }

func (s *TreeSet[T, C]) rotateRight(n *node[T]) {
	parent := n.parent
	leftChild := n.left

	n.left = leftChild.right
	if leftChild.right != nil {
		leftChild.right.parent = n
	}

	leftChild.right = n
	n.parent = leftChild

	s.replaceChild(parent, n, leftChild)
}

func (s *TreeSet[T, C]) rotateLeft(n *node[T]) {
	parent := n.parent
	rightChild := n.right

	n.right = rightChild.left
	if rightChild.left != nil {
		rightChild.left.parent = n
	}

	rightChild.left = n
	n.parent = rightChild

	s.replaceChild(parent, n, rightChild)
}

func (s *TreeSet[T, C]) replaceChild(parent, previous, next *node[T]) {
	switch {
	case parent == nil:
		s.root = next
	case parent.left == previous:
		parent.left = next
	case parent.right == previous:
		parent.right = next
	default:
		panic("node is not child of its parent")
	}

	if next != nil {
		next.parent = parent
	}
}

func (s *TreeSet[T, C]) insert(n *node[T]) bool {
	var (
		parent *node[T] = nil
		tmp    *node[T] = s.root
	)

	for tmp != nil {
		parent = tmp

		cmp := s.compare(n, tmp)
		switch {
		case cmp < 0:
			tmp = tmp.left
		case cmp > 0:
			tmp = tmp.right
		default:
			// already exists in tree
			return false
		}
	}

	n.color = red
	switch {
	case parent == nil:
		s.root = n
	case s.compare(n, parent) < 0:
		parent.left = n
	default:
		parent.right = n
	}
	n.parent = parent

	s.rebalanceInsertion(n)
	s.size++
	return true
}

func (s *TreeSet[T, C]) rebalanceInsertion(n *node[T]) {
	parent := n.parent

	// case 1: parent is nil
	// - means we are the root
	// - our color must be black
	if parent == nil {
		n.color = black
		return
	}

	// if parent is black there is nothing to do
	if parent.black() {
		return
	}

	// case 2: no grandparent
	// - implies the parent is root
	// - we must now be black
	grandparent := parent.parent
	if grandparent == nil {
		parent.color = black
		return
	}

	uncle := s.uncleOf(parent)

	switch {
	// case 3: uncle is red
	// - fix color of parent, grandparent, uncle
	// - recurse upwards as necessary
	case uncle != nil && uncle.red():
		parent.color = black
		grandparent.color = red
		uncle.color = black
		s.rebalanceInsertion(grandparent)

	case parent == grandparent.left:
		// case 4a: uncle is black
		// + node is left->right child of its grandparent
		if n == parent.right {
			s.rotateLeft(parent)
			parent = n // recolor in case 5a
		}

		// case 5a: uncle is black
		// + node is left->left child of its grandparent
		s.rotateRight(grandparent)

		// fix color of original parent and grandparent
		parent.color = black
		grandparent.color = red

		// parent is right child of grandparent
	default:
		// case 4b: uncle is black
		// + node is right->left child of its grandparent
		if n == parent.left {
			s.rotateRight(parent)
			// points to root of rotated sub tree
			parent = n // recolor in case 5b
		}

		// case 5b: uncle is black
		// + node is right->right child of its grandparent
		s.rotateLeft(grandparent)

		// fix color of original parent and grandparent
		parent.color = black
		grandparent.color = red
	}

}

func (*TreeSet[T, C]) uncleOf(n *node[T]) *node[T] {
	grandparent := n.parent
	switch {
	case grandparent.left == n:
		return grandparent.right
	case grandparent.right == n:
		return grandparent.left
	default:
		panic("bug: parent is not a child of our childs grandparent")
	}
}

func (s *TreeSet[T, C]) min(n *node[T]) *node[T] {
	if n.left == nil {
		return n
	}
	return s.min(n.left)
}

func (s *TreeSet[T, C]) max(n *node[T]) *node[T] {
	if n.right == nil {
		return n
	}
	return s.max(n.right)
}

func (s *TreeSet[T, C]) compare(a, b *node[T]) int {
	return s.comparison(a.element, b.element)
}
