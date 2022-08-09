//go:build !windows

package main

func NewListener(keySettings []ebiten.Key) func() []bool {
	return func() []bool {
		pressed := make([]bool, len(keySettings))
		for k, ek := range keySettings {
			pressed[k] = ebiten.IsKeyPressed(ek)
		}
		return pressed
	}
}
