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
	// size = 10_000
	size = 10
)

type token struct {
	id string
}

func (t *token) String() string {
	return t.id
}

func compareTokens(a, b *token) int {
	return Cmp(a.id, b.id)
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
	ts := NewTreeSet[*token, Compare[*token]](compareTokens)
	must.NotNil(t, ts)
	ts.dump()
}

func TestTreeSetFrom(t *testing.T) {
	s := shuffle(ints(10))
	ts := TreeSetFrom[int, Compare[int]](s, Cmp[int])
	must.NotEmpty(t, ts)
}

func TestTreeSet_Empty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		must.Empty(t, ts)
	})

	t.Run("not empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		ts.Insert(1)
		must.NotEmpty(t, ts)
	})
}

func TestTreeSet_Size(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		must.Size(t, 0, ts)
	})
	t.Run("one", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		ts.Insert(42)
		must.Size(t, 1, ts)
	})
	t.Run("ten", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		s := shuffle(ints(10))
		for i := 0; i < len(s); i++ {
			ts.Insert(s[i])
			must.Size(t, i+1, ts)
		}
		// insert again (all duplicates)
		s = shuffle(s)
		for i := 0; i < len(s); i++ {
			ts.Insert(s[i])
			must.Size(t, 10, ts)
		}
	})
}

// Slice

// String

// StringFunc

func TestTreeSet_Insert_token(t *testing.T) {
	ts := NewTreeSet[*token, Compare[*token]](compareTokens)

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
	cmp := Cmp[int]
	ts := NewTreeSet[int, Compare[int]](cmp)

	numbers := ints(size)
	random := shuffle(numbers)

	for _, i := range random {
		ts.Insert(i)
		invariants(t, ts, cmp)
	}

	t.Log("dump: insert int")
	t.Log(ts.dump())
}

func TestTreeSet_InsertSlice(t *testing.T) {
	cmp := Cmp[int]

	numbers := ints(size)
	random := shuffle(numbers)

	ts := NewTreeSet[int, Compare[int]](cmp)
	must.True(t, ts.InsertSlice(random))
	must.Eq(t, numbers, ts.Slice())
	must.False(t, ts.InsertSlice(numbers))
}

func TestTreeSet_Remove_int(t *testing.T) {
	cmp := Cmp[int]
	ts := NewTreeSet[int, Compare[int]](cmp)

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

	// all gone
	must.Empty(t, ts)
}

func TestTreeSet_RemoveSlice(t *testing.T) {
	cmp := Cmp[int]
	ts := NewTreeSet[int, Compare[int]](cmp)

	numbers := ints(size)
	random := shuffle(numbers)
	ts.InsertSlice(random)

	must.True(t, ts.RemoveSlice(numbers))
	must.Empty(t, ts)
}

func TestTreeSet_Contains(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		must.False(t, ts.Contains(42))
	})

	t.Run("exists", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		must.Contains[int](t, 1, ts)
		must.Contains[int](t, 2, ts)
		must.Contains[int](t, 3, ts)
		must.Contains[int](t, 4, ts)
		must.Contains[int](t, 5, ts)
	})

	t.Run("absent", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		must.NotContains[int](t, 0, ts)
		must.NotContains[int](t, 6, ts)
	})
}

func TestTreeSet_ContainsSlice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ts := NewTreeSet[int, Compare[int]](Cmp[int])
		must.False(t, ts.ContainsSlice([]int{42, 43, 44}))
	})

	t.Run("exists", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		must.True(t, ts.ContainsSlice([]int{2, 1, 3}))
		must.True(t, ts.ContainsSlice([]int{5, 4, 3, 2, 1}))
	})

	t.Run("absent", func(t *testing.T) {
		ts := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3, 4, 5}, Cmp[int])
		must.False(t, ts.ContainsSlice([]int{6, 7, 8}))
		must.False(t, ts.ContainsSlice([]int{4, 5, 6}))
	})
}

func TestTreeSet_Subset(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		t1 := NewTreeSet[int, Compare[int]](Cmp[int])
		t2 := NewTreeSet[int, Compare[int]](Cmp[int])
		must.True(t, t1.Subset(t2))
	})

	t.Run("empty full", func(t *testing.T) {
		t1 := NewTreeSet[int, Compare[int]](Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		must.False(t, t1.Subset(t2))
	})

	t.Run("full empty", func(t *testing.T) {
		t1 := NewTreeSet[int, Compare[int]](Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		must.True(t, t2.Subset(t1))
	})

	t.Run("same", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 1, 3}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{1, 2, 3}, Cmp[int])
		must.True(t, t1.Subset(t2))
		must.True(t, t2.Subset(t1))
	})

	t.Run("subset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 1, 3}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{5, 4, 1, 2, 3}, Cmp[int])
		must.False(t, t1.Subset(t2))
	})

	t.Run("superset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{5, 4, 2, 1, 3}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{5, 1, 2, 3}, Cmp[int])
		must.True(t, t1.Subset(t2))
	})
}

func TestTreeSet_Union(t *testing.T) {
	t.Run("empty empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		result := t1.Union(t2)
		must.Empty(t, result)
	})

	t.Run("empty full", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{3, 1, 2}, Cmp[int])
		result := t1.Union(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3}, result.Slice())
	})

	t.Run("full empty", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 3, 1}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]](nil, Cmp[int])
		result := t1.Union(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3}, result.Slice())
	})

	t.Run("subset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 3, 1}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{2}, Cmp[int])
		result := t1.Union(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3}, result.Slice())
	})

	t.Run("superset", func(t *testing.T) {
		t1 := TreeSetFrom[int, Compare[int]]([]int{2, 3, 1}, Cmp[int])
		t2 := TreeSetFrom[int, Compare[int]]([]int{2, 5, 1, 2, 4}, Cmp[int])
		result := t1.Union(t2)
		must.NotEmpty(t, result)
		must.Eq(t, []int{1, 2, 3, 4, 5}, result.Slice())
	})
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
func invariants[T any, C Compare[T]](t *testing.T, tree *TreeSet[T, C], cmp C) {
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

// create a copy of s and shuffle
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
