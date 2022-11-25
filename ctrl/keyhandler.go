package ctrl

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/input"
)

type KeyHandler struct {
	Handler
	Modifiers []input.Key // Handler works only when all Modifier are pressed.
	Keys      [2]input.Key
	Sounds    [2]audios.Sounder
	Volume    *float64

	holdIndex int
	// holdKey   input.Key
	active bool

	countdown    int // Require to hold for a while to move a cursor.
	MaxCountdown [2]int
	// TPS          int
}

func NewKeyHandler(tps int) (h KeyHandler) {
	for i, max := range []int64{200, 80} {
		h.MaxCountdown[i] = ToTick(max, tps)
	}
	return
}

// Update returns whether the handler has set off (triggered) or not.
// Todo: set off the stricter handler only
func (h *KeyHandler) Update() (set bool) {
	if h.countdown > 0 {
		h.countdown--
		return
	}
	for _, k := range h.Modifiers {
		if !ebiten.IsKeyPressed(k) {
			h.reset()
			return
		}
	}
	if h.holdIndex > none && !ebiten.IsKeyPressed(h.Keys[h.holdIndex]) {
		h.reset()
	}
	for i, k := range h.Keys {
		if ebiten.IsKeyPressed(k) {
			h.holdIndex = i
			break
		}
	}
	if h.holdIndex == none {
		return
	}
	[]func(){h.Decrease, h.Increase}[h.holdIndex]()
	h.Sounds[h.holdIndex].Play(*h.Volume)

	const (
		long = iota
		short
	)
	if h.active {
		h.countdown = h.MaxCountdown[short]
	} else {
		h.countdown = h.MaxCountdown[long]
	}
	h.active = true
	return true
}
func (h *KeyHandler) reset() {
	h.active = false
	h.holdIndex = none
}
