// +build !windows

package kb

import "github.com/hajimehoshi/ebiten/v2"

const numKeys = ebiten.KeyMax

func isKeyPressed(i int) bool {
	eKey := ebiten.Key(i)
	return ebiten.IsKeyPressed(eKey)
}

func getKeyCode(i int) Code {
	eKey := ebiten.Key(i)

	switch eKey {
	case ebiten.KeyArrowLeft:
		return CodeLeftArrow
	case ebiten.KeyArrowUp:
		return CodeUpArrow
	case ebiten.KeyArrowRight:
		return CodeRightArrow
	case ebiten.KeyArrowDown:
		return CodeDownArrow
	case ebiten.KeyD:
		return CodeD
	case ebiten.KeyF:
		return CodeF
	case ebiten.KeyJ:
		return CodeJ
	case ebiten.KeyK:
		return CodeK
	}

	return CodeUnknown
}
