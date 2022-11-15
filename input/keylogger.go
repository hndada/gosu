package input

type KeyLogger struct {
	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool
}

func NewKeyLogger(keySettings []Key) (k KeyLogger) {
	keyCount := len(keySettings)
	k.FetchPressed = NewListener(keySettings)
	k.LastPressed = make([]bool, keyCount)
	k.Pressed = make([]bool, keyCount)
	return
}
func (l KeyLogger) KeyAction(k int) KeyAction {
	return CurrentKeyAction(l.LastPressed[k], l.Pressed[k])
}
