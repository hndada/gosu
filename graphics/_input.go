package graphic

type InputService interface {
	CursorPosition() (x int, y int)
	InputChars() []rune
	IsKeyPressed(key Key) bool
	IsKeyJustPressed(key Key) bool
	IsKeyJustReleased(key Key) bool
	IsMouseButtonPressed(button MouseButton) bool
	IsMouseButtonJustPressed(button MouseButton) bool
	IsMouseButtonJustReleased(button MouseButton) bool
	KeyPressDuration(key Key) int
}
