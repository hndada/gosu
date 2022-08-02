package main

import "github.com/hajimehoshi/ebiten/v2"

type KeyAction int

const (
	Idle KeyAction = iota
	Hit
	Release
	Hold
)

var KeySettings = map[int][]ebiten.Key{
	4: []ebiten.Key{ebiten.KeyD, ebiten.KeyF, ebiten.KeyJ, ebiten.KeyK},
	7: []ebiten.Key{ebiten.KeyS, ebiten.KeyD, ebiten.KeyF,
		ebiten.KeySpace, ebiten.KeyJ, ebiten.KeyK, ebiten.KeyL},
}

func (s *ScenePlay) KeyAction(k int) KeyAction {
	return CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}
func CurrentKeyAction(last, now bool) KeyAction {
	switch {
	case !last && !now:
		return Idle
	case !last && now:
		return Hit
	case last && !now:
		return Release
	case last && now:
		return Hold
	default:
		panic("not reach")
	}
}
