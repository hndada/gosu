package choose

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/scene"
)

const (
	RowWidth  = 550 // 400(card) + 150(list)
	RowHeight = 75  // 50
	RowShrink = 0.15 * RowWidth
	RowCount  = 15 // int(ScreenSizeY/RowHeight) + 2
)

type Row struct {
	// thumbCh chan draws.Image
	// cardCh  chan draws.Image

	draws.Sprite // Thumbnail
	Thumb        draws.Sprite
	Card         draws.Sprite
	Mask         draws.Sprite
	First        draws.Sprite
	Second       draws.Sprite
}

var defaultThumb = draws.Image{
	Image: ebiten.NewImage(int(RowHeight), int(RowHeight))}
var defaultCard = draws.Image{
	Image: ebiten.NewImage(int(RowWidth-RowHeight), int(RowHeight))}

func NewRow(cardURL, thumbURL, first, second string) Row {
	const thumbWidth = RowHeight // Thumbnail is a square.
	const (
		px = 5
		py = 30
	)
	r := Row{}
	r.Locate(ScreenSizeX-RowWidth, ScreenSizeY/2, draws.LeftMiddle)
	{
		s := draws.NewSprite(defaultThumb)
		s.SetSize(thumbWidth, RowHeight)
		r.Thumb = s
	}
	go func() {
		i, err := draws.LoadImageFromURL(thumbURL)
		if err != nil {
			return
		}
		r.Thumb.Source = i
		// r.thumbCh <- draws.Image{Image: i}
		// close(r.thumbCh)
	}()
	{
		s := draws.NewSprite(defaultCard)
		s.SetSize(400, RowHeight)
		s.Locate(thumbWidth, 0, draws.LeftTop)
		r.Card = s
	}
	go func() {
		i, err := draws.LoadImageFromURL(cardURL)
		if err != nil {
			return
		}
		r.Card.Source = i
		// r.cardCh <- draws.Image{Image: i}
		// close(r.cardCh)
	}()
	{
		s := scene.UserSkin.BoxMask
		s.SetSize(RowWidth, RowHeight)
		s.Locate(thumbWidth, 0, draws.LeftTop)
		r.Mask = s
	}
	{
		src := draws.NewText(first, scene.Face20)
		s := draws.NewSprite(src)
		s.Locate(px+thumbWidth, py, draws.LeftTop)
		r.First = s
	}
	{
		src := draws.NewText(second, scene.Face20)
		s := draws.NewSprite(src)
		s.Locate(px+thumbWidth, py-5+RowHeight/2, draws.LeftTop)
		r.Second = s
	}
	return r
}
func (r *Row) Update() {
	// select {
	// case i := <-r.thumbCh:
	// 	r.Sprite.Source = i
	// case i := <-r.cardCh:
	// 	r.Card.Source = i
	// default:
	// }
}
func (r Row) Draw(dst draws.Image) {
	r.Thumb.Position = r.Thumb.Add(r.Position)
	r.Thumb.Draw(dst, draws.Op{})
	r.Card.Position = r.Card.Add(r.Position)
	r.Card.Draw(dst, draws.Op{})
	r.Mask.Position = r.Mask.Add(r.Position)
	r.Mask.Draw(dst, draws.Op{})
	r.First.Position = r.First.Add(r.Position)
	r.First.Draw(dst, draws.Op{})
	r.Second.Position = r.Second.Add(r.Position)
	r.Second.Draw(dst, draws.Op{})
}

type List struct {
	Rows   []Row
	Cursor ctrl.KeyHandler
	cursor int
}

func NewList(rows []Row) *List {
	l := &List{}
	l.Rows = rows
	l.Cursor = ctrl.KeyHandler{
		Handler: &ctrl.IntHandler{
			Value: &l.cursor,
			Min:   0,
			Max:   len(rows) - 1,
			Loop:  false,
		},
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{input.KeyArrowUp, input.KeyArrowDown},
		Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
		Volume:    &mode.S.VolumeSound,
	}
	return l
}

const (
	prev = iota - 1
	stay
	next
)

func (l *List) Update() (bool, int) {
	last := l.cursor
	fired := l.Cursor.Update()
	now := l.cursor
	if fired && now == last {
		switch now {
		case 0:
			return fired, prev
		case RowCount - 1:
			return fired, next
		}
	}
	return fired, stay
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
			r.Move(-RowShrink+10, 0)
		}
		r.Draw(dst)
	}
}
