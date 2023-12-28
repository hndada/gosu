//go:build !windows

package input

import "github.com/hajimehoshi/ebiten/v2"

// newFetchKeyboardState returns closure.
func newFetchKeyboardState(keys []Key) func() []bool {
	return func() []bool {
		ps := make([]bool, len(ks))
		for k, ek := range ks {
			ps[k] = ebiten.IsKeyPressed(ek)
		}
		return ps
	}
}
