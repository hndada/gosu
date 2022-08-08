package input

import (
	"time"

	"golang.org/x/sys/windows"
)

var (
	moduser32            = windows.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = moduser32.NewProc("GetAsyncKeyState")
)

type KeyHook struct {
	start       time.Time
	keyEvents   []KeyEvent
	lastPressed [256]bool
	closed      bool
}

// Todo: rename either one of 'start's
func (h *KeyHook) Listen(start time.Time) {
	h.start = start
	h.keyEvents = make([]KeyEvent, 0, 30)
	go h.scan()
}
func (h *KeyHook) Flush() []KeyEvent {
	es := h.keyEvents
	h.keyEvents = make([]KeyEvent, 0, 30)
	return es
}
func (h *KeyHook) Close() { h.closed = true }
func (h *KeyHook) scan() {
	const d = 1 * time.Millisecond
	for {
		enter := time.Now()
		t := time.Since(h.start).Milliseconds()
		for i := 0; i < int(0xFF); i++ { // Query whether keys mapped from 0x00 to 0xFF is pressed.
			code := convVirtualKeyCode(uint32(i))
			if code == 0 || code == CodeUnknown {
				continue
			}
			switch {
			case !h.lastPressed[i] && h.isKeyPressed(i):
				// fmt.Printf("%s pressed at %v ms\n", code, t)
				e := KeyEvent{
					Time:    t,
					KeyCode: code,
					Pressed: true,
				}
				h.keyEvents = append(h.keyEvents, e)
				h.lastPressed[i] = true
			case h.lastPressed[i] && !h.isKeyPressed(i):
				// fmt.Printf("%s released at %v ms\n", code, t)
				e := KeyEvent{
					Time:    t,
					KeyCode: code,
					Pressed: false,
				}
				h.keyEvents = append(h.keyEvents, e)
				h.lastPressed[i] = false
			}
		}
		remained := d - time.Since(enter) // Todo: should subtract -1?
		time.Sleep(remained)              // prevents 100% CPU usage
		if h.closed {
			return
		}
	}
}

func (h KeyHook) isKeyPressed(i int) bool {
	const (
		wasPressed = 0x0001 // deprecated: whether the key was pressed after the previous call to GetAsyncKeyState
		isPressed  = 0x8000
	)
	v, _, _ := procGetAsyncKeyState.Call(uintptr(i))
	return v&isPressed != 0
}
