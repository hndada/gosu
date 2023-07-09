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
