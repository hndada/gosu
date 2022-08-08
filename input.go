package main

import "time"

type InputMode int

var CurrentInputMode InputMode

const (
	InputModeEbiten InputMode = iota
	InputModeHook
	InputModeReplay
)

// Replay는 이미 Code->Key 까지 처리되어 있음
type KeyEvent struct {
	Time    int64
	KeyCode Code
	Pressed bool
}
type KeyListener interface {
	Listen(start time.Time)
	Flush() []KeyEvent
	Close()
}

// type KeyEvent struct {
// 	Time    int64
// 	Pressed bool
// 	Key     int // Key layout index
// }

// var KeyEvents = make([]KeyEvent, 0)

// Flush supposes Listen has already called when in hook input mode.
// func Flush() []KeyEvent {
// 	switch CurrentInputMode {
// 	case InputModeEbiten:

// 	case InputModeHook:

// 	case InputModeReplay:

// 	}
// }

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
