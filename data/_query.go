package db

import "github.com/hndada/gosu/mode"

// const (
// 	ModeAll    = iota // Todo: implement
// 	ModeTypePiano4        // 1, 2, 3, 4 Key
// 	ModeTypePiano7        // 5, 6, 7 Key
// 	ModeTypePiano8        // 8 ~ Key
// 	ModeDrum
// 	ModeJjava

// 	LastMode
// )

// Todo: should SortByName's move unit be a set of Chart?
const (
	SortByName = iota
	SortByLevel

	LastSortBy
)

var ChartBoxs = make(map[string]ChartBox) // Key is a file path.
var ChartViews = make([][]ChartBox, (LastMode-1)*(LastSortBy-1))

func ViewMode(m, sort int) int {
	var m2 int
	if m&mode.ModeTypePiano != 0 {
		switch m - mode.ModeTypePiano {
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
