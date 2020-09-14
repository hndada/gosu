package exp

import (
	"sort"
	"testing"
)

func BenchmarkAppendAtTail(b *testing.B) {
	var s = make([]int, 0)
	for i := 0; i < b.N; i++ {
		s = append(s, i)
		sort.Slice(s, func(i, j int) bool {
			return s[i] < s[j]
		})
	}
}

func BenchmarkAppendAtHead(b *testing.B) {
	var s = make([]int, 0)
	for i := 0; i < b.N; i++ {
		s = append([]int{i}, s...)
		sort.Slice(s, func(i, j int) bool {
			return s[i] > s[j]
		})
	}
}
