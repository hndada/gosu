package input

type KeyAction int

const (
	Idle KeyAction = iota
	Hit
	Release
	Hold
)

func keyAction(last, current bool) KeyAction {
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

func keyActions(last, current []bool) []KeyAction {
	a := make([]KeyAction, len(current))
	for k, p := range current {
		lp := last[k]
		a[k] = keyAction(lp, p)
	}
	return a
}
