package choose

import (
	"fmt"
	"sort"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

type Chart struct {
	BeatmapId        int
	ParentSetId      int
	DiffName         string
	FileMD5          string
	Mode             int
	BPM              float64
	AR               float64
	OD               float64
	CS               float64
	HP               float64
	TotalLength      int
	HitLength        int
	Playcount        int
	Passcount        int
	MaxCombo         int
	DifficultyRating float64
	OsuFile          string
	DownloadPath     string
}
type ChartList struct {
	*List
	Charts []*Chart
	Panel  *ChartPanel
}

func (sl ChartSetList) NewChartList() (l ChartList) {
	cs := sl.Current().ChildrenBeatmaps
	rows := make([]Row, len(cs))
	sort.Slice(rows, func(i, j int) bool {
		return cs[i].DifficultyRating < cs[j].DifficultyRating
	})
	l.List = NewList(rows)
	l.Charts = cs
	l.Panel = NewChartPanel(sl.Panel, cs[0])
	return
}
func (l *ChartList) Update() {
	if l.Panel != nil {
		l.Panel.Update()
	}
	if l.Cursor.Update() {
		// Update Background
		l.Panel = NewChartPanel(l.Panel.ChartSetPanel, l.Charts[l.cursor])
	}
}
func (l ChartList) Current() *Chart {
	if len(l.Charts) == 0 {
		return nil
	}
	return l.Charts[l.cursor]
}

func NewChartRows(css ChartSet, cs []*Chart) []Row {
	rows := make([]Row, len(cs))
	for i, c := range cs {
		var r Row
		{
			t := draws.NewText(css.Title, scene.Face16)
			r.First = draws.NewSpriteFromSource(t)
		}
		{
			lv := int(c.DifficultyRating) * 4
			src := fmt.Sprintf("(Level: %d) %s", lv, c.DiffName)
			t := draws.NewText(src, scene.Face16)
			r.Second = draws.NewSpriteFromSource(t)
		}
		rows[i] = r
	}

	return rows
}

// https://osu.ppy.sh/docs/index.html#beatmapsetcompact-covers
// slimcover: 1920x360
// cover: 900x250
// card: 400x140
// list: 150x150
const APIBeatmap = "https://assets.ppy.sh/beatmaps"
const Large = "@2x"

func (c Chart) URLDownload() string {
	return fmt.Sprintf("https://api.chimu.moe/v1/%s", c.DownloadPath)
}

// ChartPanel has own Duration and BPM.
// Todo: chart channel
type ChartPanel struct {
	*ChartSetPanel
	Duration draws.Sprite // in seconds.
	BPM      draws.Sprite

	ChartName draws.Sprite
	Level     draws.Sprite
	NoteCount draws.Sprite
}

func NewChartPanel(sp *ChartSetPanel, c *Chart) *ChartPanel {
	p := &ChartPanel{
		ChartSetPanel: sp,
	}
	{
		second := c.HitLength
		t := fmt.Sprintf("%02d:%02d", second/60, second%60)
		src := draws.NewText(t, scene.Face16)
		s := draws.NewSpriteFromSource(src)
		s.Locate(450, 0, draws.RightTop)
		p.Duration = s
	}
	{
		bpm := c.BPM
		src := draws.NewText(fmt.Sprintf("%.0f", bpm), scene.Face16)
		s := draws.NewSpriteFromSource(src)
		s.Locate(450, 50, draws.LeftTop)
		p.BPM = s
	}
	{
		src := draws.NewText(c.DiffName, scene.Face20)
		s := draws.NewSpriteFromSource(src)
		s.Locate(0, 80, draws.LeftTop)
		p.ChartName = s
	}
	{ // Todo: use gosu's own level system
		lv := int(c.DifficultyRating * 4)
		src := draws.NewText(fmt.Sprintf("Level: %2d", lv), scene.Face16)
		s := draws.NewSpriteFromSource(src)
		s.Locate(450, 100, draws.RightTop)
		p.Level = s
	}
	// Todo: NoteCount
	// Due to different logic, MaxCombo tells nothing.
	return p
}
func (p *ChartPanel) Update() {
	p.ChartSetPanel.Update()
}
func (p ChartPanel) Draw(dst draws.Image) {
	p.Sprite.Draw(dst, draws.Op{})
	p.MusicName.Draw(dst, draws.Op{})
	p.Artist.Draw(dst, draws.Op{})
	p.ChartName.Draw(dst, draws.Op{})
	p.Charter.Draw(dst, draws.Op{})
	p.UpdateDate.Draw(dst, draws.Op{})

	p.Duration.Draw(dst, draws.Op{})
	p.BPM.Draw(dst, draws.Op{})
	p.Level.Draw(dst, draws.Op{})
	// p.NoteCount.Draw(dst, draws.Op{})
}
