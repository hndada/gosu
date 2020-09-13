package mode

import (
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"time"
)

// GetAsyncKeyState: most significant bit 가 set이면 해당 키 눌린 거

type KeyState int8

const (
	KeyStateUnknown KeyState = iota - 1
	KeyStateUp
	KeyStateDown
)

type KeyEvent struct {
	Time    int64
	State   KeyState
	KeyCode types.VKCode
}

type KeyboardEventChannel struct {
	start  time.Time
	c      chan types.KeyboardEvent
	events []KeyEvent
	finish chan bool
}

func NewKeyboardEventChannel() (KeyboardEventChannel, error) {
	c := KeyboardEventChannel{
		c: make(chan types.KeyboardEvent, 100),
	}
	if err := keyboard.Install(nil, c.c); err != nil {
		return c, err
	}
	c.events = make([]KeyEvent, 0, 100)
	c.finish = make(chan bool, 1)
	return c, nil
}
func (c *KeyboardEventChannel) SetStartTime(t time.Time) { c.start = t }

func (c *KeyboardEventChannel) Close() error {
	return keyboard.Uninstall()
}
func (c *KeyboardEventChannel) Listen() {
	for {
		select {
		case k := <-c.c:
			var event KeyEvent
			event.Time = time.Since(c.start).Milliseconds()
			event.State = state(k.Message)
			event.KeyCode = k.VKCode
			c.events = append(c.events, event)
		case <-c.finish:
			return
		}
	}
}
func (c *KeyboardEventChannel) Dequeue() []KeyEvent {
	e := c.events
	c.events = make([]KeyEvent, 0, 100)
	return e
}

func state(message types.Message) KeyState {
	switch message {
	case types.WM_KEYDOWN:
		return KeyStateDown
	case types.WM_KEYUP:
		return KeyStateUp
	default:
		return KeyStateUnknown
	}
}
