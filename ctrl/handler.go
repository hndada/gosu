package ctrl

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	HandlerKeyIncrease = iota
	HandlerKeyDecrease
	HandlerKeyNext
	HandlerKeyPrev
)

var HandlerKeyTypes = []int{
	HandlerKeyIncrease,
	HandlerKeyDecrease,
	HandlerKeyNext,
	HandlerKeyPrev,
}

type Handler struct {
	Keys       []ebiten.Key
	PlaySounds []func() // SoundBytes [][]byte
	HoldKey    ebiten.Key
	Countdown  int
	Active     bool
	// Hold       int // Tick count of holding.
	// Threshold  int
}

const KeyNone = -1 // ebiten.Key(input.KeyReserved0)

func (h *Handler) Update() {
	if ebiten.IsKeyPressed(h.HoldKey) {
		h.Countdown--
		if h.Countdown < 0 {
			h.Countdown = 0
		}
	} else {
		h.Active = false
		h.Countdown = longCountdown
		h.HoldKey = KeyNone
		for _, keyType := range HandlerKeyTypes[:len(h.Keys)] {
			key := h.Keys[keyType]
			if ebiten.IsKeyPressed(key) {
				h.HoldKey = key
				h.Countdown--
			}
		}
	}
}

func (h Handler) KeyType() int {
	for i, t := range HandlerKeyTypes {
		if h.HoldKey == h.Keys[i] {
			return t
		}
	}
	return -1
}

type F64Handler struct {
	Handler
	Min    float64
	Max    float64
	Unit   float64
	Target *float64
}

// Update returns whether the handler has fired or not.
func (h *F64Handler) Update() bool {
	h.Handler.Update()
	// if h.Hold%h.Threshold != 1 {}
	if h.Countdown > 0 {
		return false
	}
	// Countdown is 0: Time to action!
	h.PlaySounds[h.KeyType()]()
	switch h.KeyType() {
	case HandlerKeyIncrease:
		*h.Target += h.Unit
		if *h.Target > h.Max {
			*h.Target = h.Max
		}
	case HandlerKeyDecrease:
		*h.Target -= h.Unit
		if *h.Target < h.Min {
			*h.Target = h.Min
		}
	case HandlerKeyNext:
	case HandlerKeyPrev:
	case -1:
		panic("invalid input at handler")
	}
	if h.Active {
		h.Countdown = shortCountdown
	} else {
		h.Countdown = longCountdown
	}
	h.Active = true
	return true
}

type IntHandler struct {
	Handler
	Min    int
	Max    int
	Unit   int
	Target *int
}

// Update returns whether the handler has fired or not.
func (h *IntHandler) Update() bool {
	h.Handler.Update()
	if h.Countdown > 0 {
		return false
	}
	h.PlaySounds[h.KeyType()]()
	switch h.KeyType() {
	case HandlerKeyIncrease:
		*h.Target += h.Unit
		if *h.Target > h.Max {
			*h.Target = h.Max
		}
	case HandlerKeyDecrease:
		*h.Target -= h.Unit
		if *h.Target < h.Min {
			*h.Target = h.Min
		}
	case HandlerKeyNext:
	case HandlerKeyPrev:
	case -1:
		panic("invalid input at handler")
	}
	if h.Active {
		h.Countdown = shortCountdown
	} else {
		h.Countdown = longCountdown
	}
	h.Active = true
	return true
}

type BoolHandler struct {
	Handler
	Target *bool
}

// Update returns whether the handler has fired or not.
func (h *BoolHandler) Update() bool {
	h.Handler.Update()
	// Bool value is updated only once regardless of hold duration.
	if h.Countdown > 0 || h.Active {
		return false
	}
	switch h.KeyType() {
	case HandlerKeyIncrease:
		*h.Target = !*h.Target
	case HandlerKeyDecrease:
	case HandlerKeyNext:
	case HandlerKeyPrev:
	case -1:
		panic("invalid input at handler")
	}
	h.Active = true
	return true
}
