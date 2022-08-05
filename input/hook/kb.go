package hook

import (
	"time"
)

// TEMP: Should I update KeyEvent logs in every update?
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
		t := time.Now()
		for i := 0; i < int(numKeys); i++ { // Query key mapped to integer `0x00` to `0xFF` if it's pressed.
			keyCode := getKeyCode(i)
			if keyCode == 0 || keyCode == CodeUnknown {
				continue
			}
			switch {
			case isKeyPressed(i) && !lastPressed[i]:
				// fmt.Printf("%s pressed at %vms\n", keyCode, time.Since(startTime).Milliseconds())
				e := KeyEvent{
					Time:    time.Since(startTime).Milliseconds(),
					KeyCode: keyCode,
					Pressed: true,
				}
				KeyEvents = append(KeyEvents, e)
				lastPressed[i] = true
			case !isKeyPressed(i) && lastPressed[i]:
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
		u := time.Now()
		wait := 1*time.Millisecond - u.Sub(t) - 1
		time.Sleep(wait) // prevents 100% CPU usage
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
