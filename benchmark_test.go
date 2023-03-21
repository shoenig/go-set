package set

import (
	"math/rand"
	"sort"
	"testing"
)

type test struct {
	size int
	name string
}

var cases []test = []test{
	{size: 10, name: "10"},
	{size: 1_000, name: "1k"},
	{size: 100_000, name: "100k"},
	{size: 1_000_000, name: "1m"},
}

func random(n int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = rand.Int()
	}
	return result
}

type hashint int

func (hi hashint) Hash() int {
	return int(hi)
}

func BenchmarkSet_Insert(b *testing.B) {
	for _, tc := range cases {
		s := From(random(tc.size))
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Insert(i)
			}
		})
	}
}

func BenchmarkHashSet_Insert(b *testing.B) {
	for _, tc := range cases {
		hs := NewHashSet[hashint, int](tc.size)
		numbers := random(tc.size)
		for i := 0; i < tc.size; i++ {
			hs.Insert(hashint(numbers[i]))
		}
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hs.Insert(hashint(i))
			}
		})
	}
}

func BenchmarkTreeSet_Insert(b *testing.B) {
	for _, tc := range cases {
		ts := TreeSetFrom[int, Compare[int]](random(tc.size), Cmp[int])
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ts.Insert(i)
			}
		})
	}
}

func BenchmarkSet_Minimum(b *testing.B) {
	for _, tc := range cases {
		s := From(random(tc.size))
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				values := s.Slice()
				sort.Ints(values)
				_ = values[0]
			}
		})
	}
}

func BenchmarkHashSet_Minimum(b *testing.B) {
	for _, tc := range cases {
		hs := NewHashSet[hashint, int](tc.size)
		numbers := random(tc.size)
		for i := 0; i < tc.size; i++ {
			hs.Insert(hashint(numbers[i]))
		}
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				values := hs.Slice()
				sort.Slice(values, func(a, b int) bool { return values[a] < values[b] })
				_ = values[0]
			}
		})
	}
}

func BenchmarkTreeSet_Minimum(b *testing.B) {
	for _, tc := range cases {
		ts := TreeSetFrom[int, Compare[int]](random(tc.size), Cmp[int])
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = ts.Min()
			}
		})
	}
}

func BenchmarkSlice_Minimum(b *testing.B) {
	for _, tc := range cases {
		slice := random(tc.size)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sort.Ints(slice)
				_ = slice[0]
			}
		})
	}
}
