package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	Box
	mode    int
	onClick func()
	pressed bool
}

const (
	ButtonModePressed = iota
	ButtonModeClicked // onClick goes called when mouse button is pressed and released.
)

func (b *Button) Hover() bool {
	return b.In(Pt(ebiten.CursorPosition()))
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
