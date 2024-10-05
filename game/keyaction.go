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
	// This is a wrong code.
	// if len(kss) == 1 {
	// 	kss = append(kss, kss[0])
	// }
	kas := make([]KeyboardAction, len(kss)-1)

	keysOld := kss[0].KeysPressed
	for i, ks := range kss[1:] {
		keysNew := ks.KeysPressed
		as := make([]int, len(keysNew))
		for k, old := range keysOld {
			new := keysNew[k]
			as[k] = KeyAction(old, new)
		}
		t := int32(ks.Time.Milliseconds())
		kas[i] = KeyboardAction{Time: t, KeysAction: as}
		keysOld = keysNew
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
