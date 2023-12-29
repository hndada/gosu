package game

import (
	"github.com/hndada/gosu/input"
)

type KeyActionType = int

// KeyAction is for handling key states conveniently.
const (
	Idle KeyActionType = iota
	Hit
	Released
	Holding
)

func KeyAction(old, now bool) KeyActionType {
	switch {
	case !old && !now:
		return Idle
	case !old && now:
		return Hit
	case old && !now:
		return Released
	case old && now:
		return Holding
	default:
		panic("not reach")
	}
}

type KeyboardAction struct {
	Time       int32
	KeysAction []KeyActionType
}

func KeyboardActions(kss []input.KeyboardState) []KeyboardAction {
	kas := make([]KeyboardAction, len(kss)-1)
	old := kss[0].PressedList
	for i, s := range kss[1:] {
		now := s.PressedList
		as := make([]int, len(now))
		for k := range old {
			as[k] = KeyAction(old[k], now[k])
		}
		t := int32(s.Time.Milliseconds())
		kas[i] = KeyboardAction{Time: t, KeysAction: as}
		old = now
	}
	return kas
}

func (ka KeyboardAction) KeysPressed() []bool {
	ps := make([]bool, len(ka.KeysAction))
	for i, a := range ka.KeysAction {
		ps[i] = a == Hit || a == Holding
	}
	return ps
}

func (ka KeyboardAction) KeysHolding() []bool {
	hs := make([]bool, len(ka.KeysAction))
	for i, a := range ka.KeysAction {
		hs[i] = a == Holding
	}
	return hs
}
