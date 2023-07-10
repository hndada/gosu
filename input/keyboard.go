package input

type Keyboard interface {
	Now() int32              // int32: Maximum duration is around 24 days.
	Fetch() []KeyboardAction // Fetch latest actions with 1ms unit.
	Poll()                   // Would do nothing on Replay.
	Pause()
	Resume()
	IsPaused() bool
	Output() []KeyboardState // Output all states.
}
type KeyboardState struct {
	Time    int32
	Pressed []bool
}
type KeyboardAction struct {
	Time   int32
	Action []KeyAction
}
