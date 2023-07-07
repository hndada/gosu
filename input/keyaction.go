package input

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

// remove soon
