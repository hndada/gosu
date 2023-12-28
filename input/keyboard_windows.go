//go:build windows

package input

import "golang.org/x/sys/windows"

const defaultPollingRate = 250.0 // Hz

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeystate
// Get'Async'KeyState would make the exectuable false-positive.
var (
	moduser32       = windows.NewLazyDLL("user32.dll")
	procGetKeyState = moduser32.NewProc("GetKeyState")
)

// Deprecated: whether the key was pressed after the previous call to GetKeyState
// const wasPressed = 0x0001
const isPressed = 0x8000

// Fetch: most passive. It just gathers data without any modification.
// Read: It is a bit more active than Fetch. It processes data.
// Listen: It is the most active. It waits for event to happen.
func newFetchKeyboardState(keys []Key) func() []bool {
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
