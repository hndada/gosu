//go:build windows

package kb

import (
	"golang.org/x/sys/windows"
)

var (
	moduser32            = windows.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = moduser32.NewProc("GetAsyncKeyState")
)

const (
	wasPressed = 0x0001 // deprecated: whether the key was pressed after the previous call to GetAsyncKeyState
	isPressed  = 0x8000
)

const numKeys = 0xFF

func getKeyCode(i int) Code {
	return convVirtualKeyCode(uint32(i))
}

func isKeyPressed(i int) bool {
	v, _, _ := procGetAsyncKeyState.Call(uintptr(i))
	return v&isPressed != 0
}
