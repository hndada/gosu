//go:build windows

package input

import (
	"golang.org/x/sys/windows"
)

var (
	moduser32            = windows.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = moduser32.NewProc("GetAsyncKeyState")
)

const (
	wasPressed = 0x0001 // Deprecated: whether the key was pressed after the previous call to GetAsyncKeyState
	isPressed  = 0x8000
)

func NewListener(keySettings []Key) func() []bool {
	vkcodes := make([]uint32, len(keySettings))
	for k, ek := range keySettings {
		vkcodes[k] = ToVirtualKey(ek)
	}
	return func() []bool {
		pressed := make([]bool, len(vkcodes))
		for k, vk := range vkcodes {
			v, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
			pressed[k] = v&isPressed != 0
		}
		return pressed
	}
}
