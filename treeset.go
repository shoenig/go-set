// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

// EqualFunc represents a type implementing Equal.
type EqualFunc[T any] interface {
	Equal(T) bool
}

// LessFunc represents a type implementing Less.
type LessFunc[T any] interface {
	Less(T) bool
}

// Sortable represents a type implementing Equal and Less.
type Sortable[T any] interface {
	EqualFunc[T]
	LessFunc[T]
}

// TreeSet provides a sorted set implementation.
//
// The underlying data-structure is a standard Red-Black Tree.
// https://en.wikipedia.org/wiki/Red–black_tree
//
// The implementation prioritizes readability over maximal optimizations.
type TreeSet[S Sortable[S]] struct {
	root *node[S]
	size int
}

func NewTreeSet[S Sortable[S]]() *TreeSet[S] {
	return &TreeSet[S]{
		root: nil,
		size: 0,
	}
}

// Insert item into t.
//
// Returns true if t was modified (item was not already in t), false otherwise.
func (t *TreeSet[S]) Insert(item S) bool {
	return t.insert(&node[S]{
		element: item,
		color:   red,
	})
}

// Min returns the smallest item in the set.
//
// Must not be called on an empty set.
func (t *TreeSet[S]) Min() S {
	if t.root == nil {
		panic("min: tree is empty")
	}
	n := t.min(t.root)
	return n.element
}

// Max returns the largest item in the set.
//
// Must not be called on an empty set.
func (t *TreeSet[S]) Max() S {
	if t.root == nil {
		panic("max: tree is empty")
	}
	n := t.max(t.root)
	return n.element
}

// Size returns the number of elements in the set.
func (t *TreeSet[S]) Size() int {
	return t.size
}

// Empty returns true if there are no elements in the set.
func (t *TreeSet[S]) Empty() bool {
	return t.Size() == 0
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

type node[S Sortable[S]] struct {
	element S
	color   color
	parent  *node[S]
	left    *node[S]
	right   *node[S]
}

func (n *node[S]) less(o *node[S]) bool {
	return n.element.Less(o.element)
}

func (n *node[S]) greater(o *node[S]) bool {
	return o.element.Less(n.element)
}

func (n *node[S]) black() bool {
	return n.color == black
}

func (n *node[S]) red() bool {
	return n.color == red
}

func (t *TreeSet[S]) locate(n *node[S], target S) *node[S] {
	if n == nil || n.element.Equal(target) {
		return n
	}

	if n.element.Less(target) {
		return t.locate(n.right, target)
	}

	return t.locate(n.left, target)
}

func (t *TreeSet[S]) rotateRight(n *node[S]) {
	parent := n.parent
	leftChild := n.left

	n.left = leftChild.right
	if leftChild.right != nil {
		leftChild.right.parent = n
	}

	leftChild.right = n
	n.parent = leftChild

	t.replaceChild(parent, n, leftChild)
}

func (t *TreeSet[S]) rotateLeft(n *node[S]) {
	parent := n.parent
	rightChild := n.right

	n.right = rightChild.left
	if rightChild.left != nil {
		rightChild.left.parent = n
	}

	rightChild.left = n
	n.parent = rightChild

	t.replaceChild(parent, n, rightChild)
}

func (t *TreeSet[S]) replaceChild(parent, previous, next *node[S]) {
	switch {
	case parent == nil:
		t.root = next
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

func (t *TreeSet[S]) insert(n *node[S]) bool {
	var (
		parent *node[S] = nil
		tmp    *node[S] = t.root
	)

	for tmp != nil {
		parent = tmp
		switch {
		case n.less(tmp):
			tmp = tmp.left
		case n.greater(tmp):
			tmp = tmp.right
		default:
			// already exists in tree
			return false
		}
	}

	n.color = red
	switch {
	case parent == nil:
		t.root = n
	case n.less(parent):
		parent.left = n
	default:
		parent.right = n
	}
	n.parent = parent

	t.rebalanceInsertion(n)
	t.size++
	return true
}

func (t *TreeSet[S]) rebalanceInsertion(n *node[S]) {
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

	uncle := t.uncleOf(parent)

	switch {
	// case 3: uncle is red
	// - fix color of parent, grandparent, uncle
	// - recurse upwards as necessary
	case uncle != nil && uncle.red():
		parent.color = black
		grandparent.color = red
		uncle.color = black
		t.rebalanceInsertion(grandparent)

	case parent == grandparent.left:
		// case 4a: uncle is black
		// + node is left->right child of its grandparent
		if n == parent.right {
			t.rotateLeft(parent)
			parent = n // recolor in case 5a
		}

		// case 5a: uncle is black
		// + node is left->left child of its grandparent
		t.rotateRight(grandparent)

		// fix color of original parent and grandparent
		parent.color = black
		grandparent.color = red

		// parent is right child of grandparent
	default:
		// case 4b: uncle is black
		// + node is right->left child of its grandparent
		if n == parent.left {
			t.rotateRight(parent)
			// points to root of rotated sub tree
			parent = n // recolor in case 5b
		}

		// case 5b: uncle is black
		// + node is right->right child of its grandparent
		t.rotateLeft(grandparent)

		// fix color of original parent and grandparent
		parent.color = black
		grandparent.color = red
	}

}

func (*TreeSet[S]) uncleOf(n *node[S]) *node[S] {
	grandparent := n.parent
	switch {
	case grandparent.left == n:
		return grandparent.right
	case grandparent.right == n:
		return grandparent.left
	default:
		panic("parent is not a child of its grandparent")
	}
}

func (t *TreeSet[S]) min(n *node[S]) *node[S] {
	if n.left == nil {
		return n
	}
	return t.min(n.left)
}

func (t *TreeSet[S]) max(n *node[S]) *node[S] {
	if n.right == nil {
		return n
	}
	return t.max(n.right)
}
