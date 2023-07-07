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

const (
	wasPressed = 0x0001 // Deprecated: whether the key was pressed after the previous call to GetKeyState
	isPressed  = 0x8000
)

func NewKeyInputListener(ks []Key) *KeyInputListener {
	const defaultPollingRate = 1 * time.Millisecond

	vkcodes := make([]uint32, len(ks))
	for k, ek := range ks {
		vkcodes[k] = ToVirtualKey(ek)
	}

	listen := func() KeyInputLog {
		now := time.Now()
		pressed := make([]bool, len(vkcodes))
		for k, vk := range vkcodes {
			v, _, _ := procGetKeyState.Call(uintptr(vk))
			pressed[k] = v&isPressed != 0
		}
		return KeyInputLog{
			Time:    now,
			Pressed: pressed,
		}
	}

	listener := &KeyInputListener{
		KeySettings: ks,
		vkcodes:     vkcodes,
		PollingRate: defaultPollingRate,
		Listen:      listen,

		Logs: make([]KeyInputLog, 0, 1000),
	}
	return listener
}
