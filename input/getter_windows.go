//go:build windows

package input

import (
	"time"

	"golang.org/x/sys/windows"
)

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeystate
// Get'Async'KeyState would make the exectuable false-positive.
var (
	moduser32       = windows.NewLazyDLL("user32.dll")
	procGetKeyState = moduser32.NewProc("GetKeyState")
)

// Deprecated: whether the key was pressed after the previous call to GetKeyState
// const wasPressed = 0x0001
const isPressed = 0x8000
const PollingRate = 1 * time.Millisecond

// newKeyStatesGetter returns closure.
func newKeyStatesGetter(keys []Key) func() []bool {
	vkcodes := make([]uint32, len(keys))
	for k, ek := range keys {
		vkcodes[k] = ToVirtualKey(ek)
	}

	return func() []bool {
		ps := make([]bool, len(vkcodes))
		for k, vk := range vkcodes {
			v, _, _ := procGetKeyState.Call(uintptr(vk))
			ps[k] = v&isPressed != 0
		}
		return ps
	}
}
