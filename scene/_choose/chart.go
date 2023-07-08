package choose

import (
	"fmt"
	"sort"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

type Chart struct {
	*ChartSet
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
	// Panel  *ChartPanel
}

// Todo: these codes are kinda messy
func (s *Scene) LoadChartList() {
	s.loading = true
	cset := s.ChartSets.Current()
	cs := cset.ChildrenBeatmaps
	sort.Slice(cs, func(i, j int) bool {
		return cs[i].DifficultyRating < cs[j].DifficultyRating
	})
	rows := make([]Row, 0, len(cs))
	for i := range cs {
		if s.levelLimit && cs[i].DifficultyRating > 4 {
			continue
		}
		rows = append(rows, cset.NewChartRow(i))
	}
	s.Charts.List = NewList(rows)
	s.Charts.Charts = cs
	// s.Charts.Panel = NewChartPanel(s.ChartSets.Panel, cs[0])
	s.Focus = FocusChart
	s.loading = false
}
func (cset ChartSet) NewChartRow(i int) (r Row) {
	card := cset.URLCover("card", "")
	thumb := cset.URLCover("list", "")
	c := cset.ChildrenBeatmaps[i]
	lv := Level(c.DifficultyRating)
	second := fmt.Sprintf("(Level %2d) %s", lv, c.DiffName)
	return NewRow(card, thumb, cset.Title, second)
}
func (l *ChartList) Update() {
	// if l.Panel != nil {
	// 	l.Panel.Update()
	// }
	// if l.Cursor.Update() {
	// 	// Update Background
	// 	l.Panel = NewChartPanel(l.Panel.ChartSetPanel, l.Charts[l.cursor])
	// }
	l.Cursor.Update()
}
func (l ChartList) Current() *Chart {
	if len(l.Charts) == 0 {
		return nil
	}
	return l.Charts[l.cursor]
}
func (l ChartList) Draw(dst draws.Image) {
	l.List.Draw(dst)
	// l.Panel.Draw(dst)
}

// list: 150x150
// card: 400x140
// cover: 900x250
// slimcover: 1920x360
// Example: https://assets.ppy.sh/beatmaps/784354/covers/slimcover@2x.jpg
// Reference: https://osu.ppy.sh/docs/index.html#beatmapsetcompact-covers
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
		s := draws.NewSprite(src)
		s.Locate(450, 0, draws.RightTop)
		p.Duration = s
	}
	{
		bpm := c.BPM
		src := draws.NewText(fmt.Sprintf("%.0f", bpm), scene.Face16)
		s := draws.NewSprite(src)
		s.Locate(450, 50, draws.LeftTop)
		p.BPM = s
	}
	{
		src := draws.NewText(c.DiffName, scene.Face20)
		s := draws.NewSprite(src)
		s.Locate(0, 80, draws.LeftTop)
		p.ChartName = s
	}
	{ // Todo: use gosu's own level system
		lv := Level(c.DifficultyRating)
		src := draws.NewText(fmt.Sprintf("Level: %2d", lv), scene.Face16)
		s := draws.NewSprite(src)
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
