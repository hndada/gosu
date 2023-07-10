package input

import "time"

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

type KeyActionLog struct {
	Time   time.Time
	Action []KeyActionType
}
