package game

import (
	"time"

	"github.com/hndada/gosu/input"
)

type KeyActionType = int

// KeyAction is for handling key states conveniently.
const (
	Idle KeyActionType = iota
	Hit
	Release
	Hold
)

func KeyAction(old, now bool) KeyActionType {
	switch {
	case !old && !now:
		return Idle
	case !old && now:
		return Hit
	case old && !now:
		return Release
	case old && now:
		return Hold
	default:
		panic("not reach")
	}
}

type KeyboardAction struct {
	Time       time.Duration
	KeyActions []KeyActionType
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
		kas[i] = KeyboardAction{Time: s.Time, KeyActions: as}
		old = now
	}
	return kas
}
