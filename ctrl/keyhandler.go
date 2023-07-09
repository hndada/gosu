package ctrl

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/input"
)

const (
	short = 80
	long  = 200
)

// Todo: Modifiers work strangely when there are plural modifiers.
type KeyHandler struct {
	Handler
	Modifiers []input.Key // Handler works only when all Modifier are pressed.
	Keys      [2]input.Key
	Sounds    [2]audios.Sounder
	Volume    *float64

	holdIndex int
	active    bool
	countdown int // User needs to hold for a while to activate.
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

	if h.active {
		h.countdown = short
	} else {
		h.countdown = long
	}
	h.active = true
	return true
}

func (h *KeyHandler) reset() {
	h.active = false
	h.holdIndex = none
}
