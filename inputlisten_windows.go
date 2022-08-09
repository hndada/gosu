//go:build windows

package main

import (
	"golang.org/x/sys/windows"
)

var (
	moduser32            = windows.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = moduser32.NewProc("GetAsyncKeyState")
)

func NewListener(keySettings []uint32) func(int64) KeysState {
	const (
		wasPressed = 0x0001 // deprecated: whether the key was pressed after the previous call to GetAsyncKeyState
		isPressed  = 0x8000
	)
	return func(now int64) KeysState {
		state := KeysState{now, make([]bool, len(keySettings))}
		for k, v := range keySettings {
			v, _, _ := procGetAsyncKeyState.Call(uintptr(v))
			state.Pressed[k] = v&isPressed != 0
		}
		return state
	}
}
