//go:build windows

package main

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
}

func NewKeyHook() *KeyHook {

}

// Todo: rename either one of 'start's
func (h *KeyHook) Listen(start) {

}
func (h *KeyHook) Flush() []KeyEvent {

}
func (h *KeyHook) listen() {
	start := time.Now()
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
	remained := 1*time.Millisecond - time.Now().Sub(start) // Todo: -1 ?
	time.Sleep(remained)                                   // prevents 100% CPU usage
}

func (h KeyHook) isKeyPressed(i int) bool {
	const (
		wasPressed = 0x0001 // deprecated: whether the key was pressed after the previous call to GetAsyncKeyState
		isPressed  = 0x8000
	)
	v, _, _ := procGetAsyncKeyState.Call(uintptr(i))
	return v&isPressed != 0
}
