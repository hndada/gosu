package choose

import (
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/defaultskin"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/scene"
)

type List struct {
	Texts  []string
	Cursor ctrl.KeyHandler
	cursor int
}

func NewList(ts []string) *List {
	l := &List{}
	l.Texts = ts
	l.Cursor = ctrl.KeyHandler{
		Handler: &ctrl.IntHandler{
			Value: &l.cursor,
			Min:   0,
			Max:   len(ts) - 1,
			Loop:  false,
		},
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{input.KeyArrowUp, input.KeyArrowDown},
		Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
		Volume:    &mode.S.VolumeSound,
	}
	return l
}

func (l *List) Update() bool { return l.Cursor.Update() }
func (l List) Current() int  { return l.cursor }

const (
	RowWidth  = 750 // 400(card) + 150(list)
	RowHeight = 50  // 75
	RowShrink = 0.15 * RowWidth
	RowCount  = int(ScreenSizeY/RowHeight) + 2
)

// Load box-mask.png from defaultskin
var boxMask = draws.LoadSprite(defaultskin.FS, "box-mask.png") // interface/

// May add extra effect to box arrangement. e.g., x -= y / 5
func (l List) Draw(dst draws.Image) {
	const (
		tx = 20
		ty = 20
	)
	for i := range l.Texts[:l.cursor] {
		if i > RowCount/2 {
			break
		}
		row := boxMask
		op := draws.Op{}
		op.GeoM.Translate(ScreenSizeX-1*RowWidth, ScreenSizeY/2-RowHeight/2-float64(i+1)*RowHeight)
		// row.Move(0, -float64(i+1)*RowHeight)
		row.Draw(dst, op)

		t := draws.NewText(l.Texts[l.cursor-i-1], draws.LoadDefaultFace(20))
		op.GeoM.Translate(tx, ty)
		t.Draw(dst, op)
	}
	for i := range l.Texts[l.cursor:] {
		if i > RowCount/2 {
			break
		}
		row := boxMask
		op := draws.Op{}
		op.GeoM.Translate(ScreenSizeX-1*RowWidth, ScreenSizeY/2-RowHeight/2+float64(i)*RowHeight)
		// row.Move(0, -float64(i+1)*RowHeight)
		row.Draw(dst, op)
		if i == 0 {
			op.GeoM.Translate(-RowShrink+10, 0)
		}
		t := draws.NewText(l.Texts[l.cursor+i], draws.LoadDefaultFace(20))
		op.GeoM.Translate(tx, ty)
		t.Draw(dst, op)
	}
}
