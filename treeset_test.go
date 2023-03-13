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
	ts.Insert(tokenB)
	ts.Insert(tokenC)
	ts.Insert(tokenD)
	ts.Insert(tokenE)
	ts.Insert(tokenF)
	ts.Insert(tokenG)
	ts.Insert(tokenH)
	fmt.Println("-- dump: --")
	fmt.Println(ts.dump())

	fmt.Println("-- slice --")
	fmt.Println(ts.Slice())

	fmt.Println("-- string --")
	fmt.Println(ts.String())

}

func TestTreeSet_Insert_int(t *testing.T) {
	ts := NewTreeSet[int, Comparison[int]](Compare[int])

	n := 20

	for i := 0; i < n; i++ {
		n := rand.Int() % n
		ts.Insert(n)
	}
	fmt.Println("-- dump: --")
	fmt.Println(ts.dump())
	fmt.Println("min:", ts.Min(), "max:", ts.Max(), "size:", ts.Size())

	fmt.Println("-- slice --")
	fmt.Println(ts.Slice())

	fmt.Println("-- string --")
	fmt.Println(ts.String())

	must.Ascending(t, ts.Slice())
	invariants(t, ts, Compare[int])
}

func (n *node[T]) String() string {
	if n.red() {
		return fmt.Sprintf("\033[1;31m%v\033[0m", n.element)
	}
	return fmt.Sprintf("%v", n.element)
}

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

func (s *TreeSet[T, C]) dump() string {
	var sb strings.Builder
	s.output("", "", s.root, &sb)
	return sb.String()
}

func invariants[T any, C Comparison[T]](t *testing.T, tree *TreeSet[T, C], cmp C) {
	// assert Slice elements are ascending
	slice := tree.Slice()
	must.AscendingFunc(t, slice, func(a, b T) bool {
		return cmp(a, b) < 1
	})

	// assert size of tree
	size := tree.Size()
	must.Eq(t, size, len(slice))

	// assert slice[0] is the minimum
	min := tree.Min()
	must.Eq(t, slice[0], min)

	// assert slice[len(slice)-1] is the maximum
	max := tree.Max()
	must.Eq(t, slice[len(slice)-1], max)
}
