package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows"
)

var (
	moduser32            = windows.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = moduser32.NewProc("GetAsyncKeyState")
)

const (
	WasPressed = 0x0001 // deprecated: whether the key was pressed after the previous call to GetAsyncKeyState
	IsPressed  = 0x8000
)

var lastPressed [256]bool

func main() {
	startTime := time.Now()
	for {
		for i := 0; i < 0xFF; i++ { // Query key mapped to integer `0x00` to `0xFF` if it's pressed.
			keyCode := convVirtualKeyCode(uint32(i))
			if keyCode == 0 || keyCode == CodeUnknown {
				continue
			}
			v, _, _ := procGetAsyncKeyState.Call(uintptr(i))
			switch {
			case v&IsPressed != 0 && !lastPressed[i]:
				fmt.Printf("%s pressed at %vms\n", keyCode, time.Since(startTime).Milliseconds())
				lastPressed[i] = true
			case v&IsPressed == 0 && lastPressed[i]:
				fmt.Printf("%s released at %vms\n", keyCode, time.Since(startTime).Milliseconds())
				lastPressed[i] = false
			}
		}
		time.Sleep(1 * time.Microsecond) // prevents 100% CPU usage
	}
}

// const (
// 	Idle    = 0x0000
// 	Release = WasPressed + Idle
// 	Press   = Idle + IsPressed
// 	Hold    = WasPressed + IsPressed
// )
