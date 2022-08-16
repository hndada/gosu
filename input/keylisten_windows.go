//go:build windows

package input

import (
	"golang.org/x/sys/windows"
)

var (
	moduser32            = windows.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = moduser32.NewProc("GetAsyncKeyState")
)

func NewListener(keySettings []Key) func() []bool {
	const (
		wasPressed = 0x0001 // deprecated: whether the key was pressed after the previous call to GetAsyncKeyState
		isPressed  = 0x8000
	)
	return func() []bool {
		pressed := make([]bool, len(keySettings))
		for k, ek := range keySettings {
			vkcode := ToVirtualKey(ek)
			v, _, _ := procGetAsyncKeyState.Call(uintptr(vkcode))
			pressed[k] = v&isPressed != 0
		}
		return pressed
	}
}
