package tools

import "sort"

type BaseStamp struct {
	time int64
}

// pointer receiver 쓰기로 했다면 통일해야함
func (s *BaseStamp) SetTime(t int64) { s.time = t }
func (s *BaseStamp) Time() int64      { return s.time }

type Stamp interface {
	SetTime(t int64)
	Time() int64
}
type Stamps []Stamp

// func: 만족하는 원소중 최소; (binary search:) 작으면 i를 당김
// stamp의 본질은 decending order
// (ascending order 기준 만족 안하는 원소 중 최대로 큰 원소의 idx)
func (s Stamps) Search(time int64) int {
	return sort.Search(len(s), func(i int) bool { return s[i].Time() <= time })
}

func (s Stamps) Sort() {
	sort.Slice(s, func(i, j int) bool { return s[i].Time() > s[j].Time() })
}
