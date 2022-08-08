package input

import "time"

// Replay has already proceeded from KeyCode (Code) to Key (int) beforehand.
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
