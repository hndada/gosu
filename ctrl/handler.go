package ctrl

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
)

type EffectPlayer struct {
	MainVolume *float64
}

var effectPlayer EffectPlayer

func (ep EffectPlayer) Play(src []byte, vol float64) {
	player := audios.Context.NewPlayerFromBytes(src)
	player.SetVolume(*ep.MainVolume * vol)
	player.Play()
}

const (
	Down = iota
	Up
)

type BoolHandler struct {
	Value  *bool
	Sounds [2][]byte
}

func (h *BoolHandler) Down() { h.swap() }
func (h *BoolHandler) Up()   { h.swap() }
func (h *BoolHandler) swap() {
	if !*h.Value {
		*h.Value = true
		effectPlayer.Play(h.Sounds[Up], 1)
	} else {
		*h.Value = false
		effectPlayer.Play(h.Sounds[Down], 1)
	}
}

type FloatHandler struct {
	Value          *float64
	Unit, Min, Max float64
	Sounds         [2][]byte
}

func (h *FloatHandler) Down() {
	*h.Value -= h.Unit
	if *h.Value < h.Min {
		*h.Value = h.Min
	}
	effectPlayer.Play(h.Sounds[Down], 1)
}
func (h *FloatHandler) Up() {
	*h.Value += h.Unit
	if *h.Value > h.Max {
		*h.Value = h.Max
	}
	effectPlayer.Play(h.Sounds[Up], 1)
}

type IntHandler struct {
	Value          *int
	Unit, Min, Max int
	Loop           bool
	Sounds         [2][]byte
}

func (h *IntHandler) Down() {
	*h.Value -= h.Unit
	if *h.Value < h.Min {
		if h.Loop {
			*h.Value = h.Max
		} else {
			*h.Value = h.Min
		}
	}
	effectPlayer.Play(h.Sounds[Down], 1)
}
func (h *IntHandler) Up() {
	*h.Value += h.Unit
	if *h.Value > h.Max {
		if h.Loop {
			*h.Value = h.Min
		} else {
			*h.Value = h.Max
		}
	}
	effectPlayer.Play(h.Sounds[Up], 1)
}

const (
	KeyIndexNone = iota - 1
	KeyIndexDown
	KeyIndexUp
)

type KeyHandler struct {
	Handler
	modifiers []ebiten.Key // Handler works only when all Modifier are pressed.
	keys      [2]ebiten.Key
	// sounds    [2][]byte
	holdIndex int // ebiten.Key
	countdown int // Require to hold for a while to move a cursor.
	active    bool
}

//	func NewKeyHandler(modifiers []ebiten.Key, keys [2]ebiten.Key, sounds [2][]byte) KeyHandler {
//		return KeyHandler{
//			Modifiers:  modifiers,
//			Keys:       keys,
//			Sounds:     sounds,
//			countdown:  0,
//			holdingKey: KeyIndexNone,
//		}
//	}
func (h *KeyHandler) SetKeys(modifiers []ebiten.Key, keys [2]ebiten.Key) {
	h.modifiers = modifiers
	h.keys = keys
	h.holdIndex = KeyIndexNone
}

// func (h *KeyHandler) SetSounds(sounds [2][]byte) {
// 	h.sounds = sounds
// }

// Update returns whether the handler has triggered or not.
func (h *KeyHandler) Update() (trigger bool) {
	if h.countdown > 0 {
		h.countdown--
		return
	}
	for _, k := range h.modifiers {
		if !ebiten.IsKeyPressed(k) {
			h.reset()
			return
		}
	}
	if k := h.keys[h.holdIndex]; !ebiten.IsKeyPressed(k) {
		h.reset()
		k2 := (h.holdIndex + 1) % 2
		if ebiten.IsKeyPressed(h.keys[k2]) {
			h.holdIndex = k2
		}
	}
	switch h.holdIndex {
	case KeyIndexNone:
		return
	case KeyIndexDown:
		h.Down()
	case KeyIndexUp:
		h.Up()
	}
	if h.active {
		h.countdown = shortCountdown
	} else {
		h.countdown = longCountdown
	}
	h.active = true
	// if h.holdIndex != KeyIndexNone {
	// 	h.PlaySound(h.sounds[h.holdIndex])
	// }
	return
}

func (h *KeyHandler) reset() {
	h.active = false
	h.countdown = 0
	h.holdIndex = KeyIndexNone
}

type Handler interface {
	Down()
	Up()
}
