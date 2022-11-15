package input

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
