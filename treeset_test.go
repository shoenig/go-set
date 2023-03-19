// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

const (
	size = 1000
)

type token struct {
	id string
}

func (t *token) String() string {
	return t.id
}

func compareTokens(a, b *token) int {
	return Compare(a.id, b.id)
}

var (
	tokenA = &token{id: "A"}
	tokenB = &token{id: "B"}
	tokenC = &token{id: "C"}
	tokenD = &token{id: "D"}
	tokenE = &token{id: "E"}
	tokenF = &token{id: "F"}
	tokenG = &token{id: "G"}
	tokenH = &token{id: "H"}
)

func TestNewTreeSet(t *testing.T) {
	ts := NewTreeSet[*token, Comparison[*token]](compareTokens)
	must.NotNil(t, ts)
	ts.dump()
}

func TestTreeSet_Insert_token(t *testing.T) {
	ts := NewTreeSet[*token, Comparison[*token]](compareTokens)

	ts.Insert(tokenA)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenB)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenC)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenD)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenE)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenF)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenG)
	invariants(t, ts, compareTokens)

	ts.Insert(tokenH)
	invariants(t, ts, compareTokens)

	t.Log("dump: insert token")
	t.Log(ts.dump())
}

func TestTreeSet_Insert_int(t *testing.T) {
	cmp := Compare[int]
	ts := NewTreeSet[int, Comparison[int]](cmp)

	numbers := ints(size)
	random := shuffle(numbers)

	for _, i := range random {
		ts.Insert(i)
		invariants(t, ts, cmp)
	}

	t.Log("dump: insert token")
	t.Log(ts.dump())
}

func TestTreeSet_Remove_int(t *testing.T) {
	cmp := Compare[int]
	ts := NewTreeSet[int, Comparison[int]](cmp)

	numbers := ints(size)
	random := shuffle(numbers)

	// insert in random order
	for _, i := range random {
		ts.Insert(i)
	}

	invariants(t, ts, cmp)

	// reshuffle
	random = shuffle(random)

	// remove every element in random order
	for _, i := range random {
		removed := ts.Remove(i)
		t.Log("dump: remove", i)
		t.Log(ts.dump())
		must.True(t, removed)
		invariants(t, ts, cmp)

	}

	// done
	must.Empty(t, ts)
}

// create a colorful representation of the element in node
func (n *node[T]) String() string {
	if n.red() {
		return fmt.Sprintf("\033[1;31m%v\033[0m", n.element)
	}
	return fmt.Sprintf("%v", n.element)
}

// output creates a colorful string representation of s
func (s *TreeSet[T, C]) output(prefix, cprefix string, n *node[T], sb *strings.Builder) {
	if n == nil {
		return
	}

	sb.WriteString(prefix)
	sb.WriteString(n.String())
	sb.WriteString("\n")

	if n.right != nil && n.left != nil {
		s.output(cprefix+"├── ", cprefix+"│   ", n.right, sb)
	} else if n.right != nil {
		s.output(cprefix+"└── ", cprefix+"    ", n.right, sb)
	}
	if n.left != nil {
		s.output(cprefix+"└── ", cprefix+"    ", n.left, sb)
	}
	if n.left == nil && n.right == nil {
		return
	}
}

// dump the output of s along with the slice string
func (s *TreeSet[T, C]) dump() string {
	var sb strings.Builder
	sb.WriteString("\ntree:\n")
	s.output("", "", s.root, &sb)
	sb.WriteString("string:")
	sb.WriteString(s.String())
	return sb.String()
}

// invariants makes basic assertions about tree
func invariants[T any, C Comparison[T]](t *testing.T, tree *TreeSet[T, C], cmp C) {
	// assert Slice elements are ascending
	slice := tree.Slice()
	must.AscendingFunc(t, slice, func(a, b T) bool {
		return cmp(a, b) < 1
	})

	// assert size of tree
	size := tree.Size()
	must.Eq(t, size, len(slice), must.Sprint("tree is wrong size"))

	if size == 0 {
		return
	}

	// assert slice[0] is the minimum
	min := tree.Min()
	must.Eq(t, slice[0], min, must.Sprint("tree contains wrong min"))

	// assert slice[len(slice)-1] is the maximum
	max := tree.Max()
	must.Eq(t, slice[len(slice)-1], max, must.Sprint("tree contains wrong max"))
}

// ints will create a []int from 1 to n
func ints(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i + 1
	}
	return s
}

// shuffle s
func shuffle(s []int) []int {
	c := make([]int, len(s))
	copy(c, s)

	n := len(c)
	for i := 0; i < n; i++ {
		swp := rand.Int31n(int32(n))
		c[i], c[swp] = c[swp], c[i]
	}
	return c
}
