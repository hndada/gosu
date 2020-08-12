package graphics

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

// mouse input 등 user와 상호작용하는 애는 struct가 절대 좌표까지 들고 있어야 함
type Button struct {
	Image     *ebiten.Image
	Rect      image.Rectangle // position
	mouseDown bool
	onPressed func(b *Button)
}

func NewButton(img *ebiten.Image, x, y int) *Button {
	b := &Button{}
	w, h := img.Size()
	b.Rect = image.Rect(x, y, x+w, y+h)
	b.Image = img
	return b
}

func (b *Button) SetOnPressed(f func(b *Button)) {
	b.onPressed = f
}

// button 수준에선 이미 좌표가 고정되어 있음
// 누르고 있는 동안에는 mouseDown만 체크해서 함수가 실행되면 안됨
// 보니까 field가 struct를 포인터로 받는 경우가 있긴 한듯. onPressed-button, scene-game
func (b *Button) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if InRect(b.Rect, x, y) {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else {
		if b.mouseDown {
			// change no field value
			if b.onPressed != nil {
				b.onPressed(b)
			}
		}
		b.mouseDown = false
	}
}

// todo: padding 추가
func (b *Button) Draw(screen *ebiten.Image) {
	p := b.Rect.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.X), float64(p.Y))
	screen.DrawImage(b.Image, op)
}

func InRect(r image.Rectangle, x, y int) bool {
	return x >= r.Min.X && x < r.Max.X && y >= r.Min.Y && y < r.Max.Y
}
