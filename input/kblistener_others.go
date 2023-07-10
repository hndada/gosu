//go:build !windows

package input

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const PollingRate = 12 * time.Millisecond

// getKeyStatesFunc returns closure.
func getKeyStatesFunc(keys []Key) func() []bool {
	return func() []bool {
		ps := make([]bool, len(ks))
		for k, ek := range ks {
			ps[k] = ebiten.IsKeyPressed(ek)
		}
		return ps
	}
}
