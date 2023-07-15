package input

type Keyboard interface {
	Fetch(int32) []KeyboardAction // Fetch latest actions with 1ms unit.
	Output() []KeyboardState      // Output all states.
	// Now() int32              // int32: Maximum duration is around 24 days.
	// Poll()
	// Pause()
	// Resume()
	// IsPaused() bool
	// Close()
}

type KeyboardState struct {
	Time    int32
	Pressed []bool
}

type KeyboardAction struct {
	Time   int32
	Action []KeyActionType
}
