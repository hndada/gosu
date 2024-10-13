package game

import (
	"sort"

	"github.com/hndada/gosu/draws"
)

const (
	OrderAscending = iota
	OrderDescending
)

const (
	GroupByMusic = iota
	GroupByLevel
)

const (
	SortByName = iota
	SortByLevel
	SortByScore
	SortByLastPlayed
	SortByUpdateDate
)

// TODO: add filter
type SearchQuery struct {
	Query      string
	GroupBy    int
	GroupOrder int
	SortBy     int
	SortOrder  int
}

type SearchResult struct {
	Charts      [][]ChartRow
	FolderNames []draws.Text
	ChartNames  [][]draws.Text
}

// Helper function to group raws by a generic key
func GroupBy[T comparable](raws []ChartRow, key func(ChartRow) T) [][]ChartRow {
	buckets := make(map[T]int)
	ccs := make([][]ChartRow, 0, len(raws))
	for _, raw := range raws {
		k := key(raw)
		idx, exists := buckets[k]
		if !exists {
			idx = len(ccs)
			buckets[k] = idx
			ccs = append(ccs, []ChartRow{})
		}
		ccs[idx] = append(ccs[idx], raw)
	}
	return ccs
}

func (db Database) Search(q SearchQuery) (r SearchResult) {
	raws := make([]ChartRow, 0, 100)
	isQueryExist := q.Query != ""
	for _, c := range db.Chart {
		if isQueryExist && !c.IsMatch(q.Query) {
			continue
		}
		raws = append(raws, c)
	}

	ccs := make([][]ChartRow, 0, 100)
	switch q.GroupBy {
	case GroupByMusic:
		ccs = GroupBy(raws, func(c ChartRow) FSFile { return c.FSFile })
		sort.Slice(ccs, func(i, j int) bool {
			return ccs[i][0].MusicName < ccs[j][0].MusicName
		})
	case GroupByLevel:
		ccs = GroupBy(raws, func(c ChartRow) int { return int(c.Level) })
		sort.Slice(ccs, func(i, j int) bool {
			return ccs[i][0].Level < ccs[j][0].Level
		})
	}

	switch q.SortBy {
	// case sortByScore:
	// case sortByLastPlayed:
	// case sortByUpdateDate:
	case SortByName:
		for i, cs := range ccs {
			sort.Slice(ccs[i], func(a, b int) bool {
				return cs[a].MusicName < cs[b].MusicName
			})
		}
	case SortByLevel:
		for i, cs := range ccs {
			sort.Slice(ccs[i], func(a, b int) bool {
				return cs[a].Level < cs[b].Level
			})
		}
	}

	t1 := make([]draws.Text, len(ccs))
	switch q.GroupBy {
	case GroupByMusic:
		for i, cs := range ccs {
			t1[i] = draws.Text{Text: cs[0].MusicString()}
		}
	case GroupByLevel:
		for i, cs := range ccs {
			t1[i] = draws.Text{Text: cs[0].LevelString()}
		}
	}

	t2 := make([][]draws.Text, len(ccs))
	for i, cs := range ccs {
		t2[i] = make([]draws.Text, len(cs))
		for j, c := range cs {
			t2[i][j] = draws.Text{Text: c.ChartString()}
		}
	}

	r.Charts = ccs
	r.FolderNames = t1
	r.ChartNames = t2
	return r
}
