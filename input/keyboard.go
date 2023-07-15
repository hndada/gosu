package input

type Keyboard interface {
	Fetch(int32) []KeyboardAction // Fetch latest actions with 1ms unit.
	Output() []KeyboardState      // Output all states.
	// Now() int32

	// Poll()
	Pause()
	Resume()
	Close()
}

type KeyboardState struct {
	Time    int32
	Pressed []bool
}

type KeyboardAction struct {
	Time   int32
	Action []KeyActionType
}
