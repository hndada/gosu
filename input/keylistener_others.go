//go:build !windows

package input

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewKeyListener(ks []Key) *KeyListener {
	const pollingRate = 12 * time.Millisecond

	listen := func() KeyPressedLog {
		now := time.Now()
		pressedList := make([]bool, len(ks))
		for k, ek := range ks {
			pressedList[k] = ebiten.IsKeyPressed(ek)
		}
		return KeyPressedLog{
			Time:        now,
			PressedList: pressedList,
		}
	}

	listener := &KeyListener{
		KeySettings: ks,
		PollingRate: pollingRate,
		Listen:      listen,

		Logs: make([]KeyPressedLog, 0, 1000),
	}
	return listener
}
