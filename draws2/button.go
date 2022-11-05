package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	subject Subject
	onClick func()
	mode    ButtonMode
	pressed bool
}
type ButtonMode int

const (
	ButtonModePressed = iota
	ButtonModeClicked // onClick goes called when mouse button is pressed and released.
)

func NewButton(s Subject, onClick func(), mode ButtonMode) Button {
	return Button{
		subject: s,
		onClick: onClick,
		mode:    mode,
	}
}

func (b *Button) Hover() bool {
	return b.subject.In(IntPt(ebiten.CursorPosition()))
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
