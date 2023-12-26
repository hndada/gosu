package base

type KeyActionType = int

// KeyAction is for handling key states conveniently.
const (
	Idle KeyActionType = iota
	Hit
	Release
	Hold
)

func KeyAction(old, new bool) KeyActionType {
	switch {
	case !old && !new:
		return Idle
	case !old && new:
		return Hit
	case old && !new:
		return Release
	case old && new:
		return Hold
	default:
		panic("not reach")
	}
}

func KeyActions(olds, news []bool) []int {
	as := make([]int, len(news))
	for k, old := range olds {
		new := news[k]
		as[k] = KeyAction(old, new)
	}
	return as
}
