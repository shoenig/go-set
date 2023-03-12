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

func TestTreeSet_Insert(t *testing.T) {
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
}

func TestTreeSet_Insert2(t *testing.T) {
	ts := NewTreeSet[*token, Comparison[*token]](compareTokens)

	n := 100

	for i := 0; i < n; i++ {
		n := rand.Int() % n
		t := &token{id: fmt.Sprintf("%02d", n)}
		ts.Insert(t)
	}
	fmt.Println("-- dump: --")
	fmt.Println(ts.dump())
	fmt.Println("min:", ts.Min(), "max:", ts.Max(), "size:", ts.Size())
}

func (n *node[S]) String() string {
	if n.red() {
		return fmt.Sprintf("\033[1;31m%v\033[0m", n.element)
	}
	return fmt.Sprintf("%v", n.element)
}

func (t *TreeSet[S, C]) append(prefix, cprefix string, n *node[S], sb *strings.Builder) {
	if n == nil {
		return
	}

	sb.WriteString(prefix)
	sb.WriteString(n.String())
	sb.WriteString("\n")

	if n.right != nil && n.left != nil {
		t.append(cprefix+"├── ", cprefix+"│   ", n.right, sb)
	} else if n.right != nil {
		t.append(cprefix+"└── ", cprefix+"    ", n.right, sb)
	}
	if n.left != nil {
		t.append(cprefix+"└── ", cprefix+"    ", n.left, sb)
	}
	if n.left == nil && n.right == nil {
		return
	}
}

func (t *TreeSet[S, C]) dump() string {
	var sb strings.Builder
	t.append("", "", t.root, &sb)
	return sb.String()
}
