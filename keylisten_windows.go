//go:build windows

package gosu

import (
	"golang.org/x/sys/windows"
)

var (
	moduser32            = windows.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = moduser32.NewProc("GetAsyncKeyState")
)

func NewListener(keySettings []uint32) func() []bool {
	const (
		wasPressed = 0x0001 // deprecated: whether the key was pressed after the previous call to GetAsyncKeyState
		isPressed  = 0x8000
	)
	return func() []bool {
		pressed := make([]bool, len(keySettings))
		for k, v := range keySettings {
			v, _, _ := procGetAsyncKeyState.Call(uintptr(v))
			pressed[k] = v&isPressed != 0
		}
		return pressed
	}
}
