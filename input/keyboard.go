package input

type Keyboard interface {
	Now() int32              // int32: Maximum duration is around 24 days.
	Fetch() []KeyboardAction // Fetch latest actions with 1ms unit.
	Poll()
	Pause()
	Resume()
	IsPaused() bool
	Output() []KeyboardState // Output all states.
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
