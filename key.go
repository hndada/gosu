package gosu

type KeyAction int

const (
	Idle KeyAction = iota
	Hit
	Release
	Hold
)

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
func (s *ScenePlay) KeyAction(k int) KeyAction {
	return CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}
