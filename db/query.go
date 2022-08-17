package db

import "github.com/hndada/gosu/mode"

const (
	ModeAll    = iota // Todo: implement
	ModePiano4        // 1, 2, 3, 4 Key
	ModePiano7        // 5, 6, 7 Key
	ModePiano8        // 8 ~ Key
	ModeDrum
	ModeJjava

	LastMode
)

// Todo: should SortByName's move unit be a set of Chart?
const (
	SortByName = iota
	SortByLevel

	LastSortBy
)

var ChartInfos = make(map[string]ChartInfo) // Key is a file path.
var ChartViews = make([][]ChartInfo, (LastMode-1)*(LastSortBy-1))

func ViewMode(m, sort int) int {
	var m2 int
	if m&mode.ModePiano != 0 {
		switch m - mode.ModePiano {
		case 1, 2, 3, 4:
			m2 = 1
		case 5, 6, 7:
			m2 = 2
		default: // More than 8
			m2 = 3
		}
	} else {
		switch {
		case m&mode.ModeDrum != 0:
			m2 = 4
		case m&mode.ModeJjava != 0:
			m2 = 4
		}
	}
	return m2*(LastSortBy-1) + sort
}
