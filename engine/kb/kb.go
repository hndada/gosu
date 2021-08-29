package kb

import (
	"time"

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

// 업데이트 때마다 마지막 index 이후 최신 log 불러오기
type KeyEvent struct {
	Time    int64 // Millisecond
	KeyCode Code
	Pressed bool
}

var (
	startTime   = time.Now()
	KeyEvents   = make([]KeyEvent, 0, 2000)
	current     int // index
	lastPressed [256]bool
	done        bool
)

func Listen() {
	done = false
	for !done {
		for i := 0; i < 0xFF; i++ { // Query key mapped to integer `0x00` to `0xFF` if it's pressed.
			keyCode := convVirtualKeyCode(uint32(i))
			if keyCode == 0 || keyCode == CodeUnknown {
				continue
			}
			v, _, _ := procGetAsyncKeyState.Call(uintptr(i))
			switch {
			case v&isPressed != 0 && !lastPressed[i]:
				// fmt.Printf("%s pressed at %vms\n", keyCode, time.Since(startTime).Milliseconds())
				e := KeyEvent{
					Time:    time.Since(startTime).Milliseconds(),
					KeyCode: keyCode,
					Pressed: true,
				}
				KeyEvents = append(KeyEvents, e)
				lastPressed[i] = true
			case v&isPressed == 0 && lastPressed[i]:
				// fmt.Printf("%s released at %vms\n", keyCode, time.Since(startTime).Milliseconds())
				e := KeyEvent{
					Time:    time.Since(startTime).Milliseconds(),
					KeyCode: keyCode,
					Pressed: false,
				}
				KeyEvents = append(KeyEvents, e)
				lastPressed[i] = false
			}
		}
		time.Sleep(1 * time.Microsecond) // prevents 100% CPU usage
	}
}
func SetTime(time time.Time) {
	startTime = time
}
func Fetch() []KeyEvent {
	//count := len(KeyEvents) - current
	if current >= len(KeyEvents) {
		return []KeyEvent{}
	}
	result := KeyEvents[current:]
	current = len(KeyEvents)
	return result
}

func Exit() {
	KeyEvents = make([]KeyEvent, 0, 2000)
	current = 0
	lastPressed = [256]bool{}
	done = true
}
