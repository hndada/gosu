package input

type KeyActionType int

const (
	Idle KeyActionType = iota
	Hit
	Release
	Hold
)

func KeyAction(last, current bool) KeyActionType {
	switch {
	case !last && !current:
		return Idle
	case !last && current:
		return Hit
	case last && !current:
		return Release
	case last && current:
		return Hold
	default:
		panic("not reach")
	}
}

func KeyActions(last, current []bool) []KeyActionType {
	a := make([]KeyActionType, len(current))
	for k, p := range current {
		lp := last[k]
		a[k] = KeyAction(lp, p)
	}
	return a
}
