package choose

import (
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/scene"
)

const (
	RowWidth  float64 = 450
	RowHeight float64 = 50
	RowShrink float64 = 0.15 * RowWidth
	RowCount  int     = int(ScreenSizeY/RowHeight) + 2
)

type Row struct {
	// base
	draws.Sprite

	// Thumb stands for thumbnail.
	// Inner does not have Thumb.
	Thumb draws.Sprite

	// First contains music name for both outer and inner.
	First draws.Sprite

	// outer: Artist
	// inner: (Level) ChartName
	Second draws.Sprite
}

func NewChartSetRow(cs *ChartSet) {
	const thumbSize = 300
	const (
		dx = 20 // Padding left.
		dy = 30 // Padding bottom.
	)
}
func NewChartRow(cs *ChartSet, c Chart) {

}
func (r Row) Draw(dst draws.Image) {
	r.Thumb.Position = r.Thumb.Add(r.Position)
	r.First.Position = r.First.Add(r.Position)
	r.Second.Position = r.Second.Add(r.Position)
	r.Sprite.Draw(dst, draws.Op{})
	r.Thumb.Draw(dst, draws.Op{})
	r.First.Draw(dst, draws.Op{})
	r.Second.Draw(dst, draws.Op{})
}

type List struct {
	Rows   []Row
	Cursor ctrl.KeyHandler
	cursor int
}

func NewList(rows []Row) (l *List) {
	return &List{
		Rows: rows,
		Cursor: ctrl.KeyHandler{
			Handler: &ctrl.IntHandler{
				Value: &l.cursor,
				Min:   0,
				Max:   len(rows) - 1,
				Loop:  true,
			},
			Modifiers: []input.Key{},
			Keys:      [2]input.Key{input.KeyArrowUp, input.KeyArrowDown},
			Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
			Volume:    &mode.S.VolumeSound,
		},
	}
}
func (l *List) Update() {
	l.Cursor.Update()
}

// May add extra effect to box arrangement. e.g., x -= y / 5
func (l List) Draw(dst draws.Image) {
	for i := range l.Rows[:l.cursor] {
		r := l.Rows[l.cursor-i-1]
		if i > RowCount/2 {
			break
		}
		r.Move(0, -float64(i+1)*RowHeight)
		r.Draw(dst)
	}
	for i, r := range l.Rows[l.cursor:] {
		if i > RowCount/2 {
			break
		}
		r.Move(0, float64(i)*RowHeight)
		if i == 0 {
			r.Move(-RowShrink, 0)
		}
		r.Draw(dst)
	}
}
