package tools

import "sort"

type BaseLog struct {
	time int64
}

func (l BaseLog) Time() int64 { return l.time }

type Log interface {
	Time() int64
}

type Logs []Log

// time series; acending order
func (ls Logs) Search(time int64) int {
	idx := sort.Search(len(ls), func(i int) bool { return ls[i].Time() >= time })
	if idx < len(ls) && ls[idx].Time() == time {
		return idx
	}
	return -1
}

// 없다면, 추가하고 Sort해야함
func (ls Logs) Sort() {
	sort.Slice(ls, func(i, j int) bool { return ls[i].Time() < ls[j].Time() })
}
