package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	Sprite
	onClick func()
	mode    ButtonMode
	pressed bool
}
type ButtonMode int

const (
	ButtonModePressed = iota
	ButtonModeClicked // onClick goes called when mouse button is pressed and released.
)

func NewButton(sprite Sprite, onClick func(), mode ButtonMode) Button {
	return Button{
		Sprite:  sprite,
		onClick: onClick,
		mode:    mode,
	}
}

func (b *Button) Hover() bool {
	return b.Sprite.In(IntVec2(ebiten.CursorPosition()))
}
func (b *Button) Update() {
	if !b.Hover() {
		b.pressed = false
	} else {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if b.mode == ButtonModePressed {
				b.onClick()
			}
			b.pressed = true
		} else {
			if b.mode == ButtonModeClicked && b.pressed {
				b.onClick()
			}
			b.pressed = false
		}
	}
}
func (b *Button) Draw(dst Image, op Op) {
	b.Sprite.Draw(dst, op)
}
