package choose

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
)

const (
	BoxWidth  float64 = 450
	BoxHeight float64 = 50
	BoxShrink float64 = 0.15 * BoxWidth
	BoxCount  int     = ScreenSizeY/BoxHeight + 2
)

type List struct {
	ChartSets      []ChartSet
	ChartSetCursor ctrl.KeyHandler
	chartSetCursor int
	Charts         []Chart
	ChartCursor    ctrl.KeyHandler
	chartCursor    int
	Inner          bool // Whether current list exposes inner list.
}

func NewCursorKeyHandler(cursor *int, len int) ctrl.KeyHandler {
	return ctrl.KeyHandler{
		Handler: &ctrl.IntHandler{
			Value: cursor,
			Min:   0,
			Max:   len - 1,
			Loop:  true,
		},
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{input.KeyArrowUp, input.KeyArrowDown},
		Sounds:    [2][]byte{SwipeSound, SwipeSound},
		Volume:    &VolumeSound,
	}
}
func NewList() {

}

// s.UpdateBackground()
func (l *List) Update() {

}
func (l *List) updateOuter() {
	sort.Slice(l.ChartSets, func(i, j int) bool {
		return l.ChartSets[i].LastUpdate < l.ChartSets[j].LastUpdate
	})
}
func (l *List) updateInner() {
	sort.Slice(l.Charts, func(i, j int) bool {
		return l.Charts[i].DifficultyRating < l.Charts[j].DifficultyRating
	})
}

type Row struct {
	draws.Sprite // base

	// Thumb stands for thumbnail.
	// Inner does not have Thumb.
	Thumb draws.Sprite

	// First contains music name for both outer and inner.
	First draws.Sprite

	// outer: Artist
	// inner: (Level) ChartName
	Second draws.Sprite
}

func NewOuterRow(cs ChartSet) (rs Row) {
	const thumbSize = 300
}
func NewInnerRows(cs ChartSet) (rs []Row) {
	rs = make([]Row, len(cs.ChildrenBeatmaps))
	for i, c := range cs.ChildrenBeatmaps {
		var r Row
		{
			t := draws.NewText(cs.Title, Face16)
			r.First = draws.NewSpriteFromSource(t)
		}
		{
			lv := int(c.DifficultyRating) * 4
			src := fmt.Sprintf("(Level: %d) %s", lv, c.DiffName)
			t := draws.NewText(src, Face16)
			r.Second = draws.NewSpriteFromSource(t)
		}
		rs[i] = r
	}
	return
}

// Currently Chart infos are not in loop.
// May add extra effect to box arrangement. e.g., x -= y / 5
func (s SceneSelect) Draw(screen draws.Image) {
	s.BackgroundDrawer.Draw(screen)
	viewport, cursor := s.Viewport()
	for i := range viewport {
		sprite := ChartItemBoxSprite
		var tx float64
		if i == cursor {
			tx -= chartInfoBoxshrink
		}
		ty := float64(i-cursor) * ChartInfoBoxHeight
		sprite.Move(tx, ty)
		sprite.Draw(screen, draws.Op{})
	}

	const (
		dx = 20 // Padding left.
		dy = 30 // Padding bottom.
	)
	for i, info := range viewport {
		sprite := ChartItemBoxSprite
		t := info.Text()
		offset := float64(i-cursor) * ChartInfoBoxHeight
		x := int(sprite.X-sprite.W()) + dx   //+ rect.Dx()
		y := int(sprite.Y-sprite.H()/2) + dy //+ rect.Dy()
		if i == cursor {
			x -= int(chartInfoBoxshrink)
		}
		text.Draw(screen.Image, t, Face12, x, y+int(offset), color.Black)
	}
}
func (l List) Draw(dst draws.Image) {
	var rs []Row
	cursor := l.chartSetCursor
	if l.Inner {
		cursor = l.chartCursor
	}
	for i := range rs[:cursor] {
		r := rs[cursor-i-1]
		if i > BoxCount/2 {
			break
		}
		r.Move(0, -float64(i+1)*BoxHeight)
		r.Draw(dst, draws.Op{})
	}
	for i, r := range rs[cursor:] {
		if i > BoxCount/2 {
			break
		}
		r.Move(0, float64(i)*BoxHeight)
		r.Draw(dst, draws.Op{})
	}
}
func (s SceneSelect) Viewport() ([]ChartInfo, int) {
	count := chartItemBoxCount
	var viewport []ChartInfo
	var cursor int
	if s.Cursor <= count/2 {
		viewport = append(viewport, s.View[0:s.Cursor]...)
		cursor = s.Cursor
	} else {
		bound := s.Cursor - count/2
		viewport = append(viewport, s.View[bound:s.Cursor]...)
		cursor = count / 2
	}
	if s.Cursor >= len(s.View)-count/2 {
		viewport = append(viewport, s.View[s.Cursor:len(s.View)]...)
	} else {
		bound := s.Cursor + count/2
		viewport = append(viewport, s.View[s.Cursor:bound]...)
	}
	return viewport, cursor
}
