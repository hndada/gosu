//go:build !windows

package input

import "github.com/hajimehoshi/ebiten/v2"

func NewListener(keySettings []Key) func() []bool {
	return func() []bool {
		pressed := make([]bool, len(keySettings))
		for k, ek := range keySettings {
			pressed[k] = ebiten.IsKeyPressed(ek)
		}
		return pressed
	}
}
