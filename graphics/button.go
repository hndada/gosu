package graphics

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

type Button struct {
	MinPt     image.Point
	Image     *ebiten.Image
	mouseDown bool
	onPressed func(b *Button)

	// Padding   image.Point // Image에서 미리 처리하고 오는게 좋을듯
}

// func NewButton(img *ebiten.Image, p image.Point) *Button {
// 	b := &Button{}
// 	b.MinPt = p
// 	b.Image = img
// 	return b
// }

// button 수준에선 이미 좌표가 고정되어 있음
// 누르고 있는 동안에는 mouseDown만 체크해서 함수가 실행되면 안됨
// 보니까 field가 struct를 포인터로 받는 경우가 있긴 한듯. onPressed-button, scene-game
func (b *Button) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p := image.Pt(ebiten.CursorPosition())
		if p.In(b.Rect()) {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else {
		if b.mouseDown {
			if b.onPressed != nil {
				b.onPressed(b)
			}
		}
		b.mouseDown = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.MinPt.X), float64(b.MinPt.Y))
	screen.DrawImage(b.Image, op)
}

func (b *Button) Rect() image.Rectangle {
	w, h := b.Image.Size()
	maxPt := b.MinPt.Add(image.Pt(w, h))
	return image.Rectangle{b.MinPt, maxPt}
}

func (b *Button) SetOnPressed(f func(b *Button)) {
	b.onPressed = f
}

// func InRect(r image.Rectangle, x, y int) bool {
// 	fmt.Println(r, x, y)
// 	return x >= r.Min.X && x < r.Max.X && y >= r.Min.Y && y < r.Max.Y
// }
